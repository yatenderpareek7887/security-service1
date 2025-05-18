package utility

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	threatentity "github.com/yatender-pareek/threat-analyzer-service/src/models/threat-model"
	"gorm.io/gorm"
)

func ProcessLogs(sqlDB *sql.DB, gormDB *gorm.DB) (int, error) {
	threats := []threatentity.Threat{}
	scanThreat := func(rows *sql.Rows, threatType, severity string) ([]threatentity.Threat, error) {
		var result []threatentity.Threat
		for rows.Next() {
			var (
				timestamp     time.Time
				userID        string
				ipAddress     string
				action        string
				fileName      *string
				databaseQuery *string
				threatType    string
			)
			if err := rows.Scan(&timestamp, &userID, &ipAddress, &action, &fileName, &databaseQuery, &threatType); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
			threat := threatentity.Threat{
				Timestamp:     timestamp,
				UserID:        userID,
				IPAddress:     ipAddress,
				Action:        action,
				FileName:      fileName,
				ThreatType:    threatType,
				Severity:      severity,
				DatabaseQuery: databaseQuery,
			}
			result = append(result, threat)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("row iteration error: %w", err)
		}
		return result, nil
	}

	// Restricted files list
	restrictedFiles := []string{
		"/secure/payroll.csv",
		"/confidential/design.pdf",
		"/db_dump.sql",
		"/public/readme.txt",
		"/logs/system.log",
	}
	fileList := "'" + strings.Join(restrictedFiles, "', '") + "'"

	// 1. Credential Stuffing
	credStuffQuery := fmt.Sprintf(`
		SELECT DISTINCT 
			l3.timestamp,
			l3.user_id,
			l3.ip_address,
			l3.action,
			l3.file_name,
			l3.database_query,
			'Credential Stuffing' AS threat_type
		FROM log_data l3
		JOIN (
			SELECT l2.user_id, l2.timestamp AS success_time
			FROM log_data l2
			WHERE l2.action = 'login_success'
			AND EXISTS (
				SELECT 1
				FROM log_data l1
				WHERE l1.user_id = l2.user_id
				AND l1.action = 'login_failed'
				AND l1.timestamp <= l2.timestamp
				AND l1.file_name IN (%s)
				GROUP BY l1.user_id
				HAVING COUNT(*) >= 3
			)
		) successful_logins
			ON l3.user_id = successful_logins.user_id
			AND l3.timestamp <= successful_logins.success_time
			AND l3.action = 'login_failed'
			AND l3.file_name IN (%s)
		ORDER BY l3.user_id, l3.timestamp;
	`, fileList, fileList)

	credRows, err := sqlDB.Query(credStuffQuery)
	if err != nil {
		return 0, fmt.Errorf("credential stuffing query failed: %w", err)
	}
	defer credRows.Close()

	credThreats, err := scanThreat(credRows, "Credential Stuffing", "High")
	if err != nil {
		return 0, fmt.Errorf("scanning credential stuffing threats failed: %w", err)
	}

	threats = append(threats, credThreats...)

	// 2. Privilege Escalation
	privEscQuery := `SELECT DISTINCT 
        lf.timestamp,
        lf.user_id,
        lf.ip_address,
        lf.action,
        lf.file_name,
        lf.database_query,
        'Privilege Escalation' AS threat_type
    FROM log_data lf
    JOIN log_data dbm
        ON dbm.user_id = lf.user_id
        AND dbm.database_query IS NOT NULL
        AND (dbm.database_query LIKE '%DELETE%' OR dbm.database_query LIKE '%INSERT%')
        AND dbm.timestamp BETWEEN lf.timestamp AND DATE_ADD(lf.timestamp, INTERVAL 5 MINUTE)
    WHERE lf.action = 'login_failed'
        AND lf.database_query IS NOT NULL
        AND (lf.database_query LIKE '%DELETE%' OR lf.database_query LIKE '%INSERT%')
    ORDER BY lf.timestamp;`

	privRows, err := sqlDB.Query(privEscQuery)
	if err != nil {
		return 0, fmt.Errorf("privilege escalation query failed: %w", err)
	}
	defer privRows.Close()

	privThreats, err := scanThreat(privRows, "Privilege Escalation", "High")
	if err != nil {
		return 0, fmt.Errorf("scanning privilege escalation threats failed: %w", err)
	}
	threats = append(threats, privThreats...)

	// 3. Account Takeover
	accTakeoverQuery := `SELECT DISTINCT 
			l3.timestamp,
			l3.user_id,
			l3.ip_address,
			l3.action,
			l3.file_name,
			l3.database_query,
			'Account Takeover' AS threat_type
		FROM log_data l3
		JOIN (
			SELECT 
				l1.user_id,
				l1.timestamp AS first_access_time
			FROM log_data l1
			JOIN log_data l2
				ON l1.user_id = l2.user_id
				AND l1.ip_address != l2.ip_address
				AND l2.timestamp BETWEEN l1.timestamp AND DATE_ADD(l1.timestamp, INTERVAL 10 MINUTE)
				AND l2.file_name IN ('/secure/payroll.csv','/db_dump.sql')
			WHERE l1.file_name IN ('/secure/payroll.csv','/db_dump.sql')
		) suspicious_access
			ON l3.user_id = suspicious_access.user_id
			AND l3.timestamp BETWEEN suspicious_access.first_access_time AND DATE_ADD(suspicious_access.first_access_time, INTERVAL 10 MINUTE)
			AND l3.file_name IN ('/confidential/design.pdf','/secure/payroll.csv','/db_dump.sql','/logs/system.log')
		ORDER BY l3.user_id, l3.timestamp;`

	accRows, err := sqlDB.Query(accTakeoverQuery)
	if err != nil {
		return 0, fmt.Errorf("account takeover query failed: %w", err)
	}
	defer accRows.Close()

	accThreats, err := scanThreat(accRows, "Account Takeover", "Medium")
	if err != nil {
		return 0, fmt.Errorf("scanning account takeover threats failed: %w", err)
	}
	threats = append(threats, accThreats...)

	// 4. Data Exfiltration
	dataExfilQuery := `SELECT DISTINCT 
			l3.timestamp,
			l3.user_id,
			l3.ip_address,
			l3.action,
			l3.file_name,
			l3.database_query,
			'Data Exfiltration' AS threat_type
		FROM log_data l3
		JOIN (
			SELECT 
				user_id, 
				MIN(timestamp) AS start_time
			FROM log_data
			WHERE action = 'file_access'
				AND file_name IN ('/confidential/design.pdf','/secure/payroll.csv','/db_dump.sql','/logs/system.log')
			GROUP BY user_id, FLOOR(UNIX_TIMESTAMP(timestamp) / 30)
			HAVING COUNT(DISTINCT file_name) >= 2
		) suspicious_access
			ON l3.user_id = suspicious_access.user_id
			AND l3.timestamp BETWEEN suspicious_access.start_time AND DATE_ADD(suspicious_access.start_time, INTERVAL 30 SECOND)
			AND l3.action = 'file_access'
			AND l3.file_name IN ('/confidential/design.pdf','/secure/payroll.csv','/db_dump.sql','/logs/system.log')
		ORDER BY l3.user_id, l3.timestamp;`

	exfilRows, err := sqlDB.Query(dataExfilQuery)
	if err != nil {
		return 0, fmt.Errorf("data exfiltration query failed: %w", err)
	}
	defer exfilRows.Close()

	exfilThreats, err := scanThreat(exfilRows, "Data Exfiltration", "High")
	if err != nil {
		return 0, fmt.Errorf("scanning data exfiltration threats failed: %w", err)
	}
	threats = append(threats, exfilThreats...)

	// 5. Insider Threat
	insiderThreatQuery := fmt.Sprintf(`
		SELECT DISTINCT 
			timestamp,
			user_id,
			ip_address,
			action,
			file_name,
			database_query,
			'Insider Threat' AS threat_type
		FROM log_data
		WHERE action = 'file_access'
			AND HOUR(timestamp) BETWEEN 2 AND 4
			AND file_name IN (%s)
		ORDER BY timestamp;
	`, fileList)

	insiderRows, err := sqlDB.Query(insiderThreatQuery)
	if err != nil {
		return 0, fmt.Errorf("insider threat query failed: %w", err)
	}
	defer insiderRows.Close()

	insiderThreats, err := scanThreat(insiderRows, "Insider Threat", "Medium")
	if err != nil {
		return 0, fmt.Errorf("scanning insider threat threats failed: %w", err)
	}
	threats = append(threats, insiderThreats...)
	for _, threat := range threats {
		if err := gormDB.Create(&threat).Error; err != nil {
			return 0, fmt.Errorf("failed to insert threat: %w", err)
		}
	}
	fmt.Printf("Inserted %d threat records into threats table\n", len(threats))
	return len(threats), nil

}

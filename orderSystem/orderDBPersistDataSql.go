package orderSystem

// client
type OrderDBPersistDataSqlClient struct {
}

func (OrderDBPersistDataSqlClient) GetInitialSql() string {
	return "call OrderClientData_Load(?,?)"
}
func (OrderDBPersistDataSqlClient) GetInsertSql() string {
	return "call OrderClientData_Insert(?,?,?,?,?,?,?)"
}
func (OrderDBPersistDataSqlClient) GetQuerySql() string {
	return "call OrderClientData_QueryByCustomId(?)"
}
func (OrderDBPersistDataSqlClient) GetUpdateSql() string {
	return "call OrderClientData_UpdateStatusByCustomId(?,?)"
}

// server
type OrderDBPersistDataSqlServer struct {
}

func (OrderDBPersistDataSqlServer) GetInitialSql() string {
	return ""
}
func (OrderDBPersistDataSqlServer) GetInsertSql() string {
	return ""
}
func (OrderDBPersistDataSqlServer) GetQuerySql() string {
	return ""
}
func (OrderDBPersistDataSqlServer) GetUpdateSql() string {
	return ""
}

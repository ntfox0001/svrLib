package database

var (
	mDatabase     *Database
	mDataResultCh chan (<-chan *DataResult)
	mQuitCh       chan interface{}
)

// func InitialDB(ip, port, user, password, database string) error {
// 	var err error = nil
// 	mDatabase, err = NewDatabase(ip, port, user, password, database, 10, 10)

// 	mDataResultCh = make(chan (<-chan *DataResult))
// 	mQuitCh = make(chan interface{}, 1)
// 	if err != nil {
// 		log.Error("database", "init", err.Error())
// 		return err
// 	}

// 	go run()

// 	return nil
// }

// func run() {
// running:
// 	for {
// 		select {
// 		case drCh := <-mDataResultCh:
// 			{
// 				go func() {
// 					t := time.NewTimer(time.Second * 3)
// 					for {
// 						select {
// 						case <-drCh:
// 							{
// 								//if dr.IsFinished() {
// 								return
// 								//}
// 							}
// 						case <-t.C:
// 							return
// 						}
// 					}
// 				}()
// 			}
// 		case <-mQuitCh:
// 			break running
// 		}
// 	}
// }

// func Query(sql string, args ...interface{}) error {

// 	op := mDatabase.CreateOperation(sql, args...)
// 	drCh, err := mDatabase.ExecOperation(op)

// 	mDataResultCh <- drCh

// 	return err
// }

// func Close() {
// 	mDatabase.Close()
// 	mQuitCh <- struct{}{}
// }

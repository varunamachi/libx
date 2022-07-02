package pg

// type JsonObject map[string]interface{}

// func (u JsonObject) Value() (driver.Value, error) {
// 	return json.Marshal(u)
// }

// func (u *JsonObject) Scan(value interface{}) error {
// 	if value == nil {
// 		return nil
// 	}
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	return json.Unmarshal(b, &u)
// }

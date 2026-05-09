package internal

func CreateVault(name string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_ = name
	return nil
}
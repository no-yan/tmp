package main

func main() {
	err := bigger()
	if err != nil {
		panic(err)
	}

	// f1, err := os.Open("sample1.txt")
	// if err != nil {
	// 	return
	// }
	//
	// f2, err := os.Open("sample2.txt")
	// if err != nil {
	// 	return
	// }
	//
	// r := io.MultiReader(f1, f2)
	// w := bufio.NewWriter(os.Stdout)
	// if _, err := io.Copy(w, r); err != nil {
	// 	log.Fatal(err)
	// }
	//
	// if err := w.Flush(); err != nil {
	// 	return
	// }
}

function test1(x) {
	if (x > 0) {
		return x;
	}

	return -x;
}

function test2(x) {
	if (x) {
		if (typeof x == "string") {
			return "string";
		}
		return "non-nullable";
	}

	return "nullable";
}

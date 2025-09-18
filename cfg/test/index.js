// function test1(x) {
// 	if (x > 0) {
// 		return x;
// 	}
//
// 	return -x;
// }
//
// function test2(x) {
// 	if (x) {
// 		if (typeof x == "string") {
// 			return "string";
// 		}
// 		return "non-nullable";
// 	}
//
// 	return "nullable";
// }

class A extends B {
	constructor() {
		if (x) {
			return;
		}
		super();
	}
}

class B {}

// class C extends B {
// 	constructor() {
// 		if (1 < 2) {
// 			super();
// 		} else {
// 		}
// 	}
// }

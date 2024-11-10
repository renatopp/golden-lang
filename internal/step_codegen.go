package internal

func CodeGen_C(pkg *Package) (string, error) {

	code := `
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
	printf("Hello, World!🥲\n");
	return 0;
}
`

	return code, nil
}

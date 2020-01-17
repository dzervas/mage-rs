#include <stdio.h>
#include "mage.h"

int main() {
	void *addr = calloc(1, 0);
	void *len = calloc(1, 0);
	char buffer[1024] = { 0 };

	/*printf("bind()\n");*/
	/*bind(0, addr, 0);*/

	printf("ffi_listen()\n");
	int fd = ffi_listen(0, 3);

	printf("ffi_accept(%d)\n", fd);
	int sock = ffi_accept(fd, addr, len);

	printf("ffi_recv(%d)\n", sock);
	ffi_recv(sock, buffer, 100, 0);
	printf("%s\n", buffer);

	printf("ffi_send(%d, world)\n", sock);
	ffi_send(sock, "world" , 6, 0);

	return 0;
}

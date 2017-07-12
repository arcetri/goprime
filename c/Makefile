SHELL= /bin/bash
CC= cc
CFLAGS= -O3 -g3 -Wall

TOPNAME= cmd
INSTALL= install

DESTDIR= /usr/local/bin

TARGETS= goprime-c

all: ${TARGETS}

goprime-c: goprime-c.c
	${CC} ${CFLAGS} goprime-c.c -D_GNU_SOURCE -lflint -lmpir -lmpfr -lm -I /usr/local/include/flint -o $@

configure:
	@echo nothing to configure

clean quick_clean quick_distclean distclean:
	rm -f goprime-c.o

clobber quick_clobber: clean
	rm -f goprime-c
	rm -rf goprime-c.dSYM

install: all
	${INSTALL} -m 0555 ${TARGETS} ${DESTDIR}

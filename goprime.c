#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <assert.h>
#include <getopt.h>
#include <stdlib.h>
#include <stdint.h>
#include <fmpz.h>
#include <errno.h>
#include <long_extras.h>
#include <stdbool.h>
#include <time.h>
#include <inttypes.h>
#include <sys/times.h>
#include <string.h>

long int debug;
char *program;
static const char *usage = "[-v] h n\n"
		"\n"
		"\t-v\tverbose mode\n"
		"\n"
		"\th\tpower of 2 multiplier (as in h*2^n-1)\n"
		"\tn\tpower of 2 (as in h*2^n-1)\n";

void dbg(int level, char const *message)
{
	if (level > debug) return;

	time_t rawtime;
	struct tm * timeinfo;

	time ( &rawtime );
	timeinfo = localtime ( &rawtime );

	char *time = asctime(timeinfo);
	time[strlen(time) - 1] = 0;

	printf("%s -> %s\n", time, message);
	fflush(stdout);

	return;
}

struct RieselNumber {
	uint64_t h;
	uint64_t n;
	fmpz_t N;	// h*2^n-1
};

struct RieselModCache {
	fmpz_t j;
	fmpz_t k;
	fmpz_t tquo;
	fmpz_t tmod;
};

void initializeRieselModCache(struct RieselModCache *C)
{
	fmpz_init(C->j);
	fmpz_init(C->k);
	fmpz_init(C->tmod);
	fmpz_init(C->tquo);
}

void clearRieselModCache(struct RieselModCache *C)
{
	fmpz_clear(C->j);
	fmpz_clear(C->k);
	fmpz_clear(C->tmod);
	fmpz_clear(C->tquo);
}

int efficientJacobi(uint64_t x, uint64_t h, uint64_t n);
uint64_t GenV1(struct RieselNumber *R);
void GenU2(fmpz_t r, struct RieselNumber *R, uint64_t v1);
void GenUN(struct RieselNumber *R, fmpz_t u2);
void rieselMod(fmpz_t a, struct RieselNumber *R, struct RieselModCache *C);

bool isPrime(struct RieselNumber *R)
{
	char *dbgMessage;

	// Check preconditions
	if (R->h < 1) {
		asprintf(&dbgMessage, "Error: expected h >= 1, but received h = %" PRId64, R->h);
		dbg(0, dbgMessage);
		return false;
	}
	if (R->n < 2) {
		asprintf(&dbgMessage, "Error: expected n >= 2, but received n = %" PRId64, R->n);
		dbg(0, dbgMessage);
		return false;
	}

	/*
	// Check if N is a small prime or a multiple of a small prime
	if check, err := screenEasyPrimes(R); err == nil && check != 0 {
		if check == 1 {
			log.Infof("N = %v is a known prime < 257", R)
			return true, nil
		}

		log.Infof("N = %v has a known factor < 257", R)
		return false, nil
	}
	*/

	// Step 1: Get a V(1) for the Riesel candidate.
	uint64_t v1 = GenV1(R);
	if (v1 == 0) { return false; }
	asprintf(&dbgMessage, "Generated V(1) = %" PRId64, v1);
	dbg(1, dbgMessage);

	// Step 2: Use the generated V(1) to generate U(2) = V(h)
	fmpz_t u;
	fmpz_init(u);
	GenU2(u, R, v1);
	if (debug == 1) {
		char *str, *message;
		fmpz_t print;
		fmpz_init(print);
		fmpz_mod_ui(print, u, 100000000);
		str = malloc(sizeof(char) * fmpz_sizeinbase(print, 10));
		fmpz_get_str(str, 10, print);
		fmpz_clear(print);
		asprintf(&message, "Generated U(2) = V(h). Last 8 digits = %s.", str);
		dbg(1, message);
		free(str);
		free(message);
	}

	// Step 3: Use the generated U(2) to generate U(n)
	GenUN(R, u);
	asprintf(&dbgMessage, "Generated U(n)");
	dbg(1, dbgMessage);

	// Step 4: Check if U(n) == 0 (mod N)
	bool result = false;
	if (fmpz_is_zero(u) == 1) {
		result = true;
	}

	fmpz_clear(u);
	free(dbgMessage);

	return result;
}

// Returns 0 if N is composite or it was not possible to generate V(1)
uint64_t GenV1(struct RieselNumber *R)
{
	char *dbgMessage;

	// Check preconditions
	if (R->h < 1) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h >= 1, but received h = %" PRId64, R->h);
		dbg(0, errorMessage);
		return 0;
	}
	if (R->n < 2) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected n >= 2, but received n = %" PRId64, R->n);
		dbg(0, errorMessage);
		return 0;
	}
	if (R->h % 2 == 0) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h mod 2 != 0, but received h = %" PRId64 "which is even", R->h);
		dbg(0, errorMessage);
		return 0;
	}

	int64_t hmod3 = R->h % 3;

	// Check if h is not a multiple of 3
	if (hmod3 != 0) {

		// Screen easy composites where 3 is a factor.
		// It is easy to show that when:
		//
		// 		(h mod 3 == 1 AND n is even) OR
		// 		(h mod 3 == 2 AND n is odd),
		//
		// then 3 is a factor.
		//
		// This relies on the observation that:
		//		2^(2k) ==  +1 (mod 3)
		//		2^(2k+1) == -1 (mod 3)
		if (((hmod3 == 1) && ((R->n & 1) == 0)) || ((hmod3 == 2) && ((R->n & 1) == 1))) {
			asprintf(&dbgMessage, "N is a multiple of 3");
			dbg(1, dbgMessage);
			return 0;
		}

		// In all these cases, we have that v(1) = 4
		return 4;
	}

	// Handle the cases when h is a multiple of 3 with the Rodseth method
	uint64_t P;
	for (P = 3; P < INT64_MAX; P++) {

		int jMinus = efficientJacobi(P - 2, R->h, R->n);
		if (jMinus == 0) {
			return 0;
		}

		if (jMinus == 1) {

			int jPlus = efficientJacobi(P + 2, R->h, R->n);
			if (jPlus == 0) {
				return 0;
			}

			if (jPlus == -1) {
				return P;
			}
		}
	}

	return 0;
}

int64_t modExp(int64_t base, uint64_t exponent, uint64_t modulus)
{
	if (base == 0 || modulus == 1) { return 0; }
	int64_t result = 1;
	base = base % modulus;

	while (exponent > 0) {
		if ((exponent & 1) == 1) {
			result = (result * base) % modulus;
		}

		exponent >>= 1;
		base = (base * base) % modulus;
	}

	return result;
}

int efficientJacobi(uint64_t x, uint64_t h, uint64_t n)
{
	bool sign = true;

	while ((x & 1) == 0) {
		x >>= 1;	// a = a / 2
		if (n == 2) { sign = !sign; }
	}

	uint64_t hModX = h % x;
	if (hModX == 0) {
		if (sign == true) {
			return 1;
		} else {
			return -1;
		}
	}

	// Jacobi(N, x) = Jacobi(((h mod x) * (2^n mod x) - 1) mod x, x).
	int64_t twoNModX = modExp(2, n, x);
	uint64_t NModX = (hModX * twoNModX - 1 + x) % x;

	// Check if x divides N (just in case)
	if ((NModX == 0) && (x != 1)) {
		char *errorMessage;
		asprintf(&errorMessage, "N has a known factor, it does not need to be tested further.");
		dbg(0, errorMessage);
		return 0;
	}

	int jNx = n_jacobi_unsigned(NModX, x);

	if (jNx == 0) {
		char *errorMessage;
		asprintf(&errorMessage, "N has a known factor, it does not need to be tested further.");
		dbg(0, errorMessage);
		return 0;
	}

	if ((x % 4) == 3) { sign = !sign; }

	// Jacobi(x, N) = Jacobi(N, x) * sign
	if (sign == true) {
		return jNx;
	} else {
		return -jNx;
	}
}

uint bitLen(uint64_t n)
{
	uint bitLen = 0;
	while(n) {
		bitLen++;
		n >>= 1;
	}

	return bitLen;
}

bool bit(int64_t n, uint index)
{
	return ((n >> index) & 1) == 1;
}

void GenU2(fmpz_t r, struct RieselNumber *R, uint64_t v1)
{
	// Check preconditions
	if (R->h < 1) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h >= 1, but received h = %" PRId64, R->h);
		dbg(0, errorMessage);
		return;
	}
	if (R->n < 2) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected n >= 2, but received n = %" PRId64, R->n);
		dbg(0, errorMessage);
		return;
	}
	if (R->h % 2 == 0) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h mod 2 != 0, but received h = %" PRId64 "which is even", R->h);
		dbg(0, errorMessage);
		return;
	}
	if (v1 < 3) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected v1 >= 3, but received v1 = %" PRId64, v1);
		dbg(0, errorMessage);
		return;
	}

	fmpz_set_ui(r, v1);

	if (R->h == 1) {
		fmpz_mod(r, r, R->N);
		return;
	}

	fmpz_t s;
	fmpz_init(s);
	fmpz_mul(s, r, r);
	fmpz_sub_ui(s, s, 2);

	uint i;
	for (i = bitLen(R->h) - 2; i > 0; i--) {

		// Starting from:
		//		r = V(x)
		//		s = V(x+1)
		if (bit(R->h, i)) {

			// If the current bit is a 1, set:
			// 		r = V(2*x+1)
			// 		s = V(2*x+2)
			//
			// These two operations are done in parallel
			fmpz_mul(r, r, s);
			fmpz_sub_ui(r, r, v1);
			fmpz_mod(r, r, R->N);

			fmpz_mul(s, s, s);
			fmpz_sub_ui(s, s, 2);
			fmpz_mod(s, s, R->N);

		} else {

			// If the current bit is a 0, set:
			// 		s = V(2*x+1)
			// 		r = V(2*x)
			//
			// These two operations are done in parallel
			fmpz_mul(s, s, r);
			fmpz_sub_ui(s, s, v1);
			fmpz_mod(s, s, R->N);

			fmpz_mul(r, r, r);
			fmpz_sub_ui(r, r, 2);
			fmpz_mod(r, r, R->N);
		}
	}

	// Since we know that h is odd, the final bit(0) is 1. Thus:
	// 		r = V(2*x+1)
	fmpz_mul(r, r, s);
	fmpz_sub_ui(r, r, v1);
	fmpz_mod(r, r, R->N);

	fmpz_clear(s);
}

void GenUN(struct RieselNumber *R, fmpz_t u)
{
	// Check preconditions
	if (R->h < 1) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h >= 1, but received h = %" PRId64, R->h);
		dbg(0, errorMessage);
		return;
	}
	if (R->n < 2) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected n >= 2, but received n = %" PRId64, R->n);
		dbg(0, errorMessage);
		return;
	}
	if (R->h % 2 == 0) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected h mod 2 != 0, but received h = %" PRId64 "which is even", R->h);
		dbg(0, errorMessage);
		return;
	}
	if (fmpz_sgn(u) < 0) {
		char *errorMessage;
		asprintf(&errorMessage, "Error: expected u > 0, but received u < 0");
		dbg(0, errorMessage);
		return;
	}

	int cmp;

	fmpz_t u_squared;
	fmpz_t j;
	fmpz_t k;
	fmpz_t j_div_h;
	fmpz_t j_mod_h;
	fmpz_t j_plus_k;

	fmpz_t print;
	fmpz_init(print);

	fmpz_init(u_squared);
	fmpz_init(j);
	fmpz_init(k);
	fmpz_init(j_div_h);
	fmpz_init(j_mod_h);
	fmpz_init(j_plus_k);

	struct tms begin, current;
	times(&begin);

	uint64_t i;
	for (i = 3; i <= R->n; i++) {

		// u = (u^2 - 2) mod N
		fmpz_mul(u_squared, u, u);
		fmpz_sub_ui(u, u_squared, 2);

		cmp = (fmpz_cmp(u, R->N));

		while (cmp >= 1) {
			fmpz_fdiv_q_2exp(j, u, R->n);
			fmpz_fdiv_r_2exp(k, u, R->n);

			if (R->h == 1) {
				fmpz_add(u, k, j);
			} else {
				fmpz_tdiv_q_ui(j_div_h, j, R->h);

				fmpz_mod_ui(j_mod_h, j, R->h);
				fmpz_mul_2exp(j_mod_h, j_mod_h, R->n);

				fmpz_add(j_plus_k, j_mod_h, k);
				fmpz_add(u, j_plus_k, j_div_h);
			}

			cmp = (fmpz_cmp(u, R->N));
		}

		if (cmp == 0) { fmpz_zero(u); }

		if (debug == 1 && i % 1000 == 0) {
			char *str, *dbgMessage;

			fmpz_mod_ui(print, u, 100000000);

			str = malloc(sizeof(char) * fmpz_sizeinbase(print, 10));
			fmpz_get_str(str, 10, print);
			times(&current);

			asprintf(&dbgMessage, "Reached U(%ld). Last 8 digits = %s. Utime = %.2f. Stime = %.2f.",
				 i, str, (float) (current.tms_utime - begin.tms_utime) / 100.0,
				 (float) (current.tms_stime - begin.tms_stime) / 100.0);


			dbg(1, dbgMessage);

			free(str);
			free(dbgMessage);
		}
	}

	fmpz_clear(print);

	fmpz_clear(u_squared);
	fmpz_clear(j);
	fmpz_clear(k);
	fmpz_clear(j_div_h);
	fmpz_clear(j_mod_h);
	fmpz_clear(j_plus_k);
}

int main(int argc, char *argv[])
{
	char *h_arg, *n_arg;
	uint64_t h, n;
	int c;

	debug = 0;

	// Parse args
	program = argv[0];
	while ((c = getopt(argc, argv, "v")) != -1) {
		switch (c) {
			case 'v':
				debug = 1;
				break;
			default:
				fprintf(stderr, "usage: %s %s", program, usage);
				exit(2);
		}
	}

	argv += (optind - 1);
	argc -= (optind - 1);
	if (argc != 3) {
		fprintf(stderr, "usage: %s %s", program, usage);
		exit(3);
	}

	h_arg = argv[1];
	errno = 0;
	h = strtoul(h_arg, NULL, 0);
	if (errno != 0 || h <= 0) {
		fprintf(stderr, "%s: FATAL: h must an integer > 0\n", program);
		fprintf(stderr, "usage: %s %s", program, usage);
		exit(4);
	}

	n_arg = argv[2];
	errno = 0;
	n = strtoul(n_arg, NULL, 0);
	if (errno != 0 || n <= 0) {
		fprintf(stderr, "%s: FATAL: n must an integer > 0\n", program);
		fprintf(stderr, "usage: %s %s", program, usage);
		exit(5);
	}

	// Force h to become odd
	if (h % 2 == 0) {
		while (h % 2 == 0 && h > 0) {
			h >>= 1;
			++n;
		}

		if (h <= 0) {
			fprintf(stderr, "%s: FATAL: new equivalent h: %lu <= 0\n", program, h);
			exit(6);
		}
	}

	struct RieselNumber *R = malloc(sizeof(struct RieselNumber));
	R->h = h;
	R->n = n;

	fmpz_init(R->N);
	fmpz_set_ui(R->N, R->h);
	fmpz_mul_2exp(R->N, R->N, R->n);
	fmpz_sub_ui(R->N, R->N, 1);

	printf("%d\n", isPrime(R));
}

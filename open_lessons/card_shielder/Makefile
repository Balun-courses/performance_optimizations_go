CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
PPROFDIR=${CURDIR}/pprof
PACKAGE=${CURDIR}/cmd/card_shielder

# ---------------------------------------- PROFILING ----------------------------------------
.PHONY: cpuprof
cpuprof:
	(PPROF_TMPDIR=${PPROFDIR} go tool pprof -http :8081 -seconds 20 http://127.0.0.1:8080/debug/pprof/profile)

.PHONE: memprof
memprof:
	(PPROF_TMPDIR=${PPROFDIR} go tool pprof -http :8081 http://127.0.0.1:8080/debug/pprof/heap)
# ---------------------------------------- PROFILING ----------------------------------------

# ---------------------------------------- SERVER ----------------------------------------
.PHONY: server
server: build
	${BINDIR}/server
# ---------------------------------------- SERVER ----------------------------------------

# ---------------------------------------- TANKER ----------------------------------------
.PHONY: tanker
tanker:
	go build -o ${BINDIR}/tanker ${CURDIR}/cmd/tanker
	${BINDIR}/tanker
# ---------------------------------------- TANKER ----------------------------------------

# ---------------------------------------- PROFILER ----------------------------------------
.PHONY: cpu_profile
cpu_profile: build
	go tool pprof ${BINDIR}/card_shielder http://127.0.0.1:8801/debug/pprof/profile

.PHONY: mem_profile
mem_profile: build
	go tool pprof ${BINDIR}/card_shielder http://127.0.0.1:8801/debug/pprof/heap
# ---------------------------------------- PROFILER ----------------------------------------

# ---------------------------------------- BUILD ----------------------------------------
.PHONY: build
build: bind-deps
	go mod tidy
	go build -o ${BINDIR}/server -tags=optimized_iterative ${PACKAGE}
# ---------------------------------------- BUILD ----------------------------------------


# ---------------------------------------- RUN ----------------------------------------
.PHONY: run
run: build
	chmod +x ${BINDIR}/server
	${BINDIR}/server
# ---------------------------------------- RUN ----------------------------------------


# ---------------------------------------- DEPENDENCIES ----------------------------------------
.PHONY: bind-deps
bind-deps: bindir

bindir:
	rm -rf ${BINDIR}
	mkdir -p ${BINDIR}
# ---------------------------------------- DEPENDENCIES ----------------------------------------

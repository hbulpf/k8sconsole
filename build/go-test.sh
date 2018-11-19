# Script parameters.
COVERAGE_REPORT_FILE=${1}
MAIN_PKG_NAME=${2}

# Install packages that are dependencies of the test. Do not run the test. Improves performance.
go test -i ${MAIN_PKG_NAME}/...

# Create coverage report file.
set -e
[ -e ${COVERAGE_REPORT_FILE} ] && rm ${COVERAGE_REPORT_FILE}
mkdir -p "$(dirname ${COVERAGE_REPORT_FILE})" && touch ${COVERAGE_REPORT_FILE}

# Run coverage tests of all project packages (without -race parameter to improve performance).
for PKG in $(go list ${MAIN_PKG_NAME}/... | grep -v vendor); do
    go test -coverprofile=profile.out -covermode=atomic ${PKG}
    if [ -f profile.out ]; then
        cat profile.out >> ${COVERAGE_REPORT_FILE}
        rm profile.out
    fi
done
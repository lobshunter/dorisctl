#!/usr/bin/bash

set -eu

# var from go template
DEPLOY_DIR={{.DeployDir}}
INSTALL_JAVA={{.InstallJava}}
IS_MASTER={{.IsMaster}}
export DORIS_HOME=$DEPLOY_DIR
export PID_DIR="${DORIS_HOME}"
if $IS_MASTER; then
    MASTER_HELPER="" # master doesn't need --helper option
else
    MASTER_HELPER="--helper {{.MasterHelper}}"
fi

declare JAVA
if [[ $INSTALL_JAVA ]]; then
    JAVA=${DEPLOY_DIR}/jdk/bin/java
    export JAVA_HOME=${DEPLOY_DIR}/jdk
else
    JAVA="$(which java)" # NOTE: only JDK 9+ is supported
fi

if [[ ! -x "${JAVA}" ]]; then
    echo "The JAVA_HOME environment variable is not defined correctly"
    echo "This environment variable is needed to run this program"
    echo "NB: JAVA_HOME should point to a JDK not a JRE"
    exit 1
fi

LOG_DIR="${DORIS_HOME}/log"

# const var
JAVA_OPTS="-Xmx1024m"
export DORIS_LOG_TO_STDERR=1

DORIS_FE_JAR=""
CLASSPATH=${CLASSPATH:-""}
for f in "${DORIS_HOME}/lib"/*.jar; do
    if [[ "${f}" == *"doris-fe.jar" ]]; then
        DORIS_FE_JAR="${f}"
        continue
    fi
    CLASSPATH="${f}:${CLASSPATH}"
done

# make sure the doris-fe.jar is at first order, so that some classed
# with same qualified name can be loaded priority from doris-fe.jar
CLASSPATH="${DORIS_FE_JAR}:${CLASSPATH}"
export CLASSPATH="${CLASSPATH}:${DORIS_HOME}/lib:${DORIS_HOME}/conf"

mkdir -p "${LOG_DIR}"
${JAVA} ${JAVA_OPTS} -XX:-OmitStackTraceInFastThrow -XX:OnOutOfMemoryError="kill -9 %p" org.apache.doris.DorisFE ${MASTER_HELPER}

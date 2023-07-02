#!/usr/bin/env bash

set -eu

# var from go template
DEPLOY_DIR={{.DeployDir}}
INSTALL_JAVA={{.InstallJava}}
export DORIS_HOME=$DEPLOY_DIR
export PID_DIR="${DORIS_HOME}"
# const var
LOG_DIR="${DORIS_HOME}/log"
export ODBCSYSINI="${DORIS_HOME}/conf"
#filter known leak for lsan.
export LSAN_OPTIONS="suppressions=${DORIS_HOME}/conf/asan_suppr.conf"
## set asan and ubsan env to generate core file
export ASAN_OPTIONS=symbolize=1:abort_on_error=1:disable_coredump=0:unmap_shadow_on_exit=1:detect_container_overflow=0
export UBSAN_OPTIONS=print_stacktrace=1
# see https://github.com/apache/doris/blob/master/docs/zh-CN/community/developer-guide/debug-tool.md#jemalloc-heap-profile
export JEMALLOC_CONF="percpu_arena:percpu,background_thread:true,metadata_thp:auto,muzzy_decay_ms:30000,dirty_decay_ms:30000,oversize_threshold:0,lg_tcache_max:16,prof_prefix:jeprof.out"

export UDF_RUNTIME_DIR="${DORIS_HOME}/lib/udf-runtime"

mkdir -p "${LOG_DIR}"
mkdir -p "${UDF_RUNTIME_DIR}"
rm -f "${UDF_RUNTIME_DIR}"/*

if [[ $INSTALL_JAVA ]]; then
    # be needs libjvm.so to boot
    export JAVA_HOME=${DEPLOY_DIR}/jdk
fi

DORIS_LOG_TO_STDERR=1 ${DORIS_HOME}/lib/doris_be

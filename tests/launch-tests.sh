#!/usr/bin/env sh

PATH=$PATH:"$PWD/node_modules/.bin":"$PWD/..":"$PWD"
ROOT_PATH=$PWD

setUp() {
    cd $ROOT_PATH/scenarios
    pkill taskmasterd

    true
}

tearDown() {
    echo $PWD
    pkill taskmasterd

    git restore .

    true
}

testInfinite() {
    cd infinite
    rm -f taskmasterd.log

    ./test.sh

    assertTrue $?

    git diff --exit-code . > /dev/null

    assertTrue "Files should have not been modified but were" $?
}

testHotReloadTotalNewConfig() {
    cd hot-reload-total-new-config
    rm -f taskmasterd.log

    ./test.sh

    assertTrue $?

    git diff --exit-code . > /dev/null

    assertFalse "Files should have been modified but were not" $?
}

testHotReloadUpdateProgramConfig() {
    cd hot-reload-update-program-config
    rm -f taskmasterd.log

    ./test.sh

    assertTrue $?

    git diff --exit-code . > /dev/null

    assertFalse "Files should have been modified but were not" $?
}

testNotFoundCommand() {
    cd not-found-command
    rm -f taskmasterd.log

    ./test.sh

    assertTrue $?

    git diff --exit-code . > /dev/null

    assertTrue "Files should have not been modified but were" $?

}

testCreateProgram() {
    cd create-program
    rm -f taskmasterd.log

    ./test.sh

    assertTrue $?

    git diff --exit-code . > /dev/null

    assertFalse "Files should have been modified but were not" $?
}

. ./vendor/shunit2/shunit2

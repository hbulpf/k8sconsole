/**
 * @fileOverview Help function that spawns go binary process.
 */
import child from 'child_process';
import lodash from 'lodash';
import q from 'q';
import semver from 'semver';

import config from "./config";

// Add base dir to the gopath so that local imports work.
const sourceGopath = `${config.paths.backendTmp}:${config.paths.backendVendor}`;
// Add project's required go tools to the PATH
const devPath = `${process.env.PATH}:${config.paths.goTools}/bin`;

/**
 * The env needed for execution of any go command
 */
const env = lodash.merge(process.env, {GOPATH: sourceGopath, PATH: devPath});

/**
 * Spawns a Go process after making sure all Go prerequisites are
 * present.
 *
 * Backend source files must be packaged with 'package-backend'
 * task before running this command.
 * @param args - Arguments of the Go command
 * @param doneFn - Callback
 * @param {!Object<string, string>} envOverride - optional environment variables override map.
 */
export default function goCommand(args, doneFn, envOverride) {
    checkPrerequisites()
        .then(() => spawnGoPrcocess(args, envOverride))
        .then(doneFn)
        .fail((error) => doneFn(error));
}

/**
 * Minimum required Go version
 */
const miniGoVersion = '1.8.3';

/**
 * Spawns a goimports process after making sure all Go prerequisites are present.
 * @param args - Arguments of the Go command
 * @param doneFn - Callback
 * @param {!Object<string, string>} envOverride - optional environment variables override map.
 */
export function goimportsCommand(args, doneFn, envOverride) {
    checkPrerequisites()
        .then(() => spawnGoimportsProcess(args, envOverride))
        .then(doneFn)
        .fail((error) => doneFn(error));
}

/**
 * Checks if all prerequisites for go-command exec are present.
 * @return {Q.Promise<Error>}
 */
function checkPrerequisites() {
    return checkGo().then(checkGoVersion).then(checkGovendor);
}

/**
 * Checks if go is on the PATH prior to a go command exec, promises an error otherwise.
 * @return {Q.Promise<Error>}
 */
function checkGo() {
    let deferred = q.defer();
    child.exec(
        'which go',
        {
            env: env,
        },
        function (error, stdout, stderror) {
            if (error || stderror || !stdout) {
                deferred.reject(new Error(
                    'Go is not on the path. Please pass the PATH variable when you run ' +
                    'the gulp task with "PATH=$PATH" or install go if you have not yet.'
                ));
            }
            deferred.resolve();
        });
    return deferred.promise;
}

/**
 * Checks if go version fulfills the minimum version prerequisite, promises an error otherwise.
 * @return {Q.Promise<Error>}
 */
function checkGoVersion() {
    let deferred = q.defer();
    child.exec(
        'go version',
        {
            env: env,
        },
        function (error, stdout) {
            let match = /go version devel/.exec(stdout.toString());
            if (match && match.length > 0) {
                // If running a development version of Go we assume the version to be
                // good enough.
                deferred.resolve();
                return;
            }

            match = /[\d.]+/.exec(stdout.toString());   // match version number
            if (match && match.length < 1) {
                deferred.reject(new Error('Go version not found.'));
                return;
            }

            //semver requires patch number, so if go version doesn't have patch, we add '.0'
            let currentGoVersion = match[0];
            if (currentGoVersion.split('.').length === 2) {
                currentGoVersion = `${currentGoVersion}.0`;
            }
            if (semver.lt(currentGoVersion, miniGoVersion)) {
                deferred.reject(new Error(
                    `The current go version '${currentGoVersion}' is older than ` +
                    `the minimum required version '${miniGoVersion}'. ` +
                    `Please upgrade your go version!`));
                return;
            }
            deferred.resolve();
        });
    return deferred.promise;
}

/**
 * Checks if govendor is on the PATH prior to a go command execution,  promises an error otherwise.
 * @return {Q.Promise<Error>}
 */
function checkGovendor() {
    let deferred = q.defer();
    child.exec(
        'which govendor',
        {
            env: env,
        },
        function (error, stdout, stderror) {
            if (error || stderror || !stdout) {
                deferred.reject(new Error(
                    'Govendor is not on the path. ' +
                    'Please run "npm install" in the base directory of the project.'));
                return;
            }
            deferred.resolve();
        });
    return deferred.promise;
}

/**
 * Spawns a process, promises an error if the go command process fails.
 *
 * @param processName
 * @param args
 * @param envOverride
 * @return {Q.Promise<any>}
 */
function spawnProcess(processName, args, envOverride) {
    let deferred = q.defer();
    let envLocal = lodash.merge(env, envOverride);
    let goTask = child.spawn(processName, args, {
        env: envLocal,
        stdio: 'inherit',
    });

    goTask.on('exit', function (code) {
        if (code !== 0) {
            deferred.reject(new Error(`Go command error, code: ${code}`));
            return;
        }
        deferred.resolve();
    });
    return deferred.promise;
}

/**
 * Spawns a go process, promises an error if the go command process fails.
 *
 * @param args
 * @param envOverride
 * @return {Q.Promise<any>}
 */
function spawnGoPrcocess(args, envOverride) {
    return spawnProcess('go', args, envOverride);
}

/**
 * Spawns goimports process, promises an error if the go command process fails.
 * @param args
 * @param envOverride
 * @return {Q.Promise<any>}
 */
function spawnGoimportsProcess(args, envOverride) {
    return spawnProcess('goimports', args, envOverride);
}

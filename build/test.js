import childProcess from 'child_process';
import gulp from 'gulp';
import gulpProtractor from 'gulp-protractor';
import karma from 'karma';
import path from 'path';

import {browserSyncInstance} from './serve';
import config from './config';
import goCommand from './gocommand';

/**
 *
 * @param {boolean} singleRun
 * @param {function(?Error=)} doneFn
 */
function runFrontendUnitTests(singleRun, doneFn) {
    let localConfig = {
        configFile: config.paths.karmaConf,
        singleRun: singleRun,
        autoWatch: !singleRun,
    };

    let server = new karma.Server(localConfig, function (failCount) {
        doneFn(failCount ? new Error("Failed" + failCount + " tests.") : undefined);
    });
    server.start();
}

function runProtractorTests(doneFn) {
    gulp.src(path.join(config.paths.integrationTest, '**/*.js'))
        .pipe(gulpProtractor.protractor({
            configFile: config.paths.protractorConf,
        }))
        .on('error',  function(err) {
            // Close browser sync server to prevent the process from hanging.
            browserSyncInstance.exit();
            // Kill backend server and cluster, if running.
            gulp.start('kill-backend');
            doneFn(err);
        })
        .on('end', function() {
            // Close browser sync server to prevent the process from hanging.
            browserSyncInstance.exit();
            // Kill backend server and cluster, if running.
            gulp.start('kill-backend');
            doneFn();
        });
}

/**
 * Runs once all unit tests of the application.
 */
gulp.task('test', ['frontend-test', 'backend-test-with-coverage']);

/**
 * Runs once all unit tests of the frontend application.
 */
gulp.task('frontend-test', ['set-test-node-env'], function(doneFn) {
    runFrontendUnitTests(true, doneFn);
});

/**
 * Runs once all unit tests of the backend application.
 */
gulp.task('backend-test', ['package-backend'], function(doneFn) {
    goCommand(config.backend.testCommandArgs, doneFn);
});

/**
 * Runs once all unit tests of the backend application with coverage report.
 */
gulp.task('backend-test-with-coverage', ['package-backend'], function(doneFn) {
    let testProcess = childProcess.execFile(
        config.paths.goTestScript, [config.paths.coverageBackend, config.backend.mainPackageName]);

    testProcess.stdout.pipe(process.stdout);
    testProcess.stderr.pipe(process.stderr);

    testProcess.on('close', (code) => {
        if (code !== 0) {
            return doneFn(new Error(`Process exited with code: ${code}`));
        }

        return doneFn();
    });
});

/**
 * Runs all unit tests of the application. Watches for changes in the source files to rerun
 * the tests.
 */
gulp.task('test:watch', ['frontend-test:watch', 'backend-test:watch']);

/**
 * Runs frontend backend application tests. Watches for changes in the source files to rerun
 * the tests.
 */
gulp.task('frontend-test:watch', ['set-test-node-env'], function(doneFn) {
    runFrontendUnitTests(false, doneFn);
});

/**
 * Runs backend application tests. Watches for changes in the source files to rerun
 * the tests.
 */
gulp.task('backend-test:watch', ['backend-test'], function() {
    gulp.watch([path.join(config.paths.backendSrc, '**/*.go')], ['backend-test']);
});


/**
 * Runs application integration tests. Uses development version of the application.
 */
gulp.task('integration-test', ['serve:nowatch', 'webdriver-update'], runProtractorTests);


/**
 * Runs application integration tests. Uses production version of the application.
 */
gulp.task('integration-test:prod', ['serve:prod', 'webdriver-update'], runProtractorTests);


/**
 * Downloads and updates webdriver. Required to keep it up to date.
 */
gulp.task('webdriver-update', gulpProtractor.webdriver_update);


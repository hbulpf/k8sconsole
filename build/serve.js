/**
 * @fileOverview Gulp tasks that serve the application
 */
import browserSync from 'browser-sync';
import browserSyncSpa from 'browser-sync-spa';
import child from 'child_process';
import gulp from 'gulp';
import path from 'path';

import config from "./config";

/**
 * Browser sync instance that serves the application
 */
export const browserSyncInstance = browserSync.create();

/**
 * Currently running backend process obejct. Null if backend is not running
 */
let runningBackendProcess = null;

/**
 * Builds array of args for backend process based on env vars and dev/prod mode.
 * @param {string} mode
 * @return {!Array<string>}
 * TODO(wzt3309) The backend args has not be designed
 */
function getBackendArgs(mode) {
    let args = [];

    if (mode === config.build.prod) {
        console.log('backend in prod mode');
    }

    if (mode === config.build.dev) {
        console.log('backend in dev mode');
    }

    return args;
}

/**
 * Initializes BrowserSync tool. Files are served from baseDir directory list
 *
 * @param {!Array<string>|string} baseDir
 */
function browserSyncInit(baseDir) {
    browserSyncInstance.use(browserSyncSpa({
        // Only for Angular apps
        selector: '[ng-app]',
    }));


    let conf = {
        browser: [],        // Needed so that the browser does not auto-launch.
        startPath: '/',
        directory: false,   // Disable directory listings
        server: {
            baseDir: baseDir,
            routes: {
                '/node_modules': config.paths.nodeModules,
            },
        },
        notify: false,
    };

    browserSyncInstance.init(conf);
}

/**
 * Serves the application in dev mode
 */
function serveDevelopmentMode() {
    browserSyncInit([
        config.paths.serve,
        config.paths.app,
    ]);
}

/**
 * Serves the application in development mode. Watches for changes in the source files to rebuild
 * development artifacts.
 */
gulp.task('serve', ['spawn-backend', 'watch'], serveDevelopmentMode);
gulp.task('serve:no-backend', ['watch'], serveDevelopmentMode);

/**
 * Serves the application in development mode.
 */
gulp.task('serve:nowatch', ['spawn-backend', 'index'], serveDevelopmentMode);

/**
 * Serves the application in prod mode.
 */
gulp.task('serve:prod', ['spawn-backend:prod']);

/**
 * Spawns new backend application process and finishes the task immediately. Previously spawned
 * backend process is killed beforehand, if any.
 *
 * The frontend pages are served by BrowserSync.
 */
gulp.task('spawn-backend', ['backend', 'kill-backend'], function () {
    runningBackendProcess = child.spawn(
        path.join(config.paths.serve, config.backend.binaryName), getBackendArgs(config.build.dev),
        {stdio: 'inherit', cwd: config.paths.serve});

    runningBackendProcess.on("exit", function () {
        // Mark there is not backend process running anymore.
        runningBackendProcess = null;
    });
});

/**
 * Spawns new backend application process and finishes the task immediately. Previously spawned
 * backend process is killed beforehand, if any.
 *
 * In production the backend does serve the frontend pages as well.
 */
gulp.task('spawn-backend:prod', ['build-frontend', 'backend:prod', 'kill-backend'], function() {
    runningBackendProcess = child.spawn(
        path.join(config.paths.dist, config.backend.binaryName), getBackendArgs(config.build.prod),
        {stdio: 'inherit', cwd: config.paths.dist});

    runningBackendProcess.on('exit', function() {
        // Mark that there is no backend process running anymore.
        runningBackendProcess = null;
    });
});

/**
 * Kills running backend process (if any)
 */
gulp.task('kill-backend', function (doneFn) {
    if (runningBackendProcess) {
        runningBackendProcess.on('exit', function () {
            // Mark that there is no backend process running anymore.
            runningBackendProcess = null;
            // Finish the task only when the backend is actually killed.
            doneFn();
        });

        runningBackendProcess.kill();
    } else {
        doneFn();
    }
});


/**
 * Watch for changes in source files and run gulp tasks to rebuild them.
 */
gulp.task('watch', ['index', 'angular-templates'], function () {
    gulp.watch([path.join(config.paths.frontendSrc, 'index.html'), 'package.json'], ['index']);

    gulp.watch(
        [
            path.join(config.paths.frontendSrc, '**/*.scss'),
        ],
        function (event) {
            if (event.type === 'changed') {
                // If is file changed, rebuild style files
                gulp.start('styles');
            } else {
                // If is file new/del, everything has to rebuilt
                gulp.start('index');
            }
    });

    gulp.watch([path.join(config.paths.frontendSrc, '**/*.js')], ['scripts-watch']);
    gulp.watch(path.join(config.paths.frontendSrc, '**/*.html'), ['angular-templates']);
    gulp.watch(path.join(config.paths.backendSrc, '**/*.go'), ['spawn-backend']);
});

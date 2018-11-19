import child from 'child_process';
import gulp from 'gulp';
import path from 'path';

import conf from './config';

function spawnDockerProcess(args, doneFn) {
    let dockerTask = child.spawn('docker', args);

    dockerTask.on('exit', function (code) {
        if (code === 0) {
            doneFn();
        } else {
            doneFn(new Error('Docker command error, code: ' + code));
        }
    });

    dockerTask.stdout.on('data', function (data) {
        console.log('' + data);
    });

    dockerTask.stderr.on('data', function (data) {
        console.error('' + data);
    });
}

/**
 * Creates Docker image for the application. The image is tagged with the image name configuration
 * constant.
 *
 * In order to run the image on a Kubernates cluster, it has to be deployed to a registry.
 */
gulp.task('docker-image', ['docker-file'], function(doneFn) {
    spawnDockerProcess([
        'build',
        // Remove intermediate containers after a successful build.
        '--rm=true',
        '--tag', conf.deploy.imageName,
        // dist is the build context
        conf.paths.dist,
    ], doneFn);
});

/**
 * Processes the Docker file and places it in the dist folder for building.
 */
gulp.task('docker-file', function() {
    return gulp.src(path.join(conf.paths.deploySrc, 'Dockerfile'))
        .pipe(gulp.dest(conf.paths.dist));
});

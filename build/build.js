/**
 * @fileoverview Gulp tasks for building the project.
 */
import del from 'del';
import gulp from 'gulp';
import config from './config';

/**
 * Builds production package for current architecture and places it in the dist directory.
 */
gulp.task('build', ['backend:prod']);

/**
 * Builds production packages for all supported architectures and places them in the dist directory.
 */
gulp.task('build:cross', ['backend:prod:cross']);

/**
 * Cleans all build artifacts.
 */
gulp.task('clean', ['clean-dist'], function () {
    return del([config.paths.tmp, config.paths.coverage]);
});

/**
 * Cleans all build artifacts in the dist/ folder.
 */
gulp.task('clean-dist', function() {
    return del([config.paths.distRoot, config.paths.distPre]);
});
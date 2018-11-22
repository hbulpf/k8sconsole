import path from 'path';

/**
 * Load the i18n and l10n configuration. Used when dashboard is built in production.
 */
let localization = require('../i18n/locale_conf.json');

/**
 * project base path
 */
const basePath = path.join(__dirname, '../');

/**
 *  Architecture configuration.
 */
const arch = {
    /**
     * Default architecture that the project is compiled to.
     * Used for local dev and test
     */
    default: 'amd64',

    // /**
    //  * List of all
    //  */
    // list: ['amd64', 'arm', 'arm64', 'ppc64le', 's390x'],
};

const os = {
    /**
     * Default os platform
     */
    default: 'linux',
    /**
     * List of all
     */
    list: ['linux', 'darwin', 'windows'],
};

const version = {
    /**
     * Current release version
     */
    release: 'v0.0.1',
    /**
     * Version name of the head release of the project
     */
    head: 'head',
};

export default {
    recordVersionExpression:
        `-X github.com/wzt3309/k8sconsole/src/app/backend/client.Version=${version.release}`,
    /**
     * Absolute paths to known dirs, e.g., to src dir
     */
    paths: {
        app: path.join(basePath, 'src/app'),
        assets: path.join(basePath, 'src/app/assets'),
        backendSrc: path.join(basePath, 'src/app/backend'),
        backendTmp: path.join(basePath, '.tmp/backend'),
        backendTmpSrc: path.join(
            basePath, '.tmp/backend/src/github.com/wzt3309/k8sconsole/src/app/backend'),
        backendTmpSrcVendor: path.join(
            basePath, '.tmp/backend/src/github.com/wzt3309/k8sconsole/vendor'),
        backendVendor: path.join(basePath, 'vendor'),
        base: basePath,
        bowerComponents: path.join(basePath, 'bower_components'),
        build: path.join(basePath, 'build'),
        coverage: path.join(basePath, 'coverage'),
        coverageBackend: path.join(basePath, 'coverage/go.txt'),
        coverageFrontend: path.join(basePath, 'coverage/lcov/lcov.info'),
        deploySrc: path.join(basePath, 'src/app/deploy'),
        dist: path.join(basePath, 'dist', os.default),
        distCross: os.list.map((os) => path.join(basePath, 'dist', os)),
        distPre: path.join(basePath, '.tmp/dist'),
        distPublic: path.join(basePath, 'dist', os.default, 'public'),
        distPublicCross: os.list.map((os) => path.join(basePath, 'dist', os, 'public')),
        distRoot: path.join(basePath, 'dist'),
        externs: path.join(basePath, 'src/app/externs'),
        frontendSrc: path.join(basePath, 'src/app/frontend'),
        frontendTest: path.join(basePath, 'src/test/frontend'),
        goTools: path.join(basePath, '.tools/go'),
        goTestScript: path.join(basePath, 'build/go-test.sh'),
        i18nProd: path.join(basePath, '.tmp/i18n'),
        integrationTest: path.join(basePath, 'src/test/integration'),
        karmaConf: path.join(basePath, 'build/karma.conf.js'),
        messagesForExtraction: path.join(basePath, '.tmp/messages_for_extraction'),
        nodeModules: path.join(basePath, 'node_modules'),
        partials: path.join(basePath, '.tmp/partials'),
        prodTmp: path.join(basePath, '.tmp/prod'),
        protractorConf: path.join(basePath, 'build/protractor.conf.js'),
        serve: path.join(basePath, '.tmp/serve'),
        src: path.join(basePath, 'src'),
        tmp: path.join(basePath, '.tmp'),
        xtbgenerator: path.join(basePath, '.tools/xtbgenerator/bin/XtbGenerator.jar'),
    },

    /**
     * The name of the Angular module
     */
    frontend: {
        /**
         * The name of Angular module, i.e., the module that bootstrap the application
         */
        rootModuleName: 'k8sconsole',
    },

    backend: {
        /**
         * The name of the backend binary.
         */
        binaryName: 'k8sconsole',
        /**
         * Name of the main backend package that is used in go build command.
         */
        mainPackageName: 'github.com/wzt3309/k8sconsole/src/app/backend',

        testCommandArgs:
            [
                'test',
                'github.com/wzt3309/k8sconsole/src/app/backend/...',
            ],
    },

    build: {
        prod: 'production',
        test: 'test',
        dev: 'development',
    },

    deploy: {
        imageName: 'wzt3309/k8sconsole',
    },

    arch: arch,
    os: os,

    /**
     * Configuration for i18n & l10n.
     */
    translations: localization.translations.map((translation) => {
        return {path: path.join(basePath, 'i18n', translation.file), key: translation.key};
    }),

};

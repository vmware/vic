/*
 Copyright 2017 VMware, Inc. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

var webpackConfig = require('./webpack.test.js');

module.exports = function (config) {
    config.set({
        basePath: '..',
        frameworks: ['jasmine', 'source-map-support'],
        plugins: [
            require('karma-source-map-support'),
            require('karma-webpack'),
            require('karma-jasmine'),
            require('karma-chrome-launcher'),
            require('karma-phantomjs-launcher'),
            require('karma-coverage-istanbul-reporter')
        ],
        customLaunchers: {
            // chrome setup for travis CI using chromium
            Chrome_travis_ci: {
                base: 'Chrome',
                flags: ['--no-sandbox']
            }
        },
        files: [
            'config/karma-entry.ts'
        ],
        preprocessors: {
            'config/karma-entry.ts': ['webpack']
        },
        webpack: webpackConfig,
        webpackMiddleware: {
            stats: 'errors-only'
        },
        reporters: ['progress', 'coverage-istanbul'],
        coverageIstanbulReporter: {
            reports: ['text-summary', 'html'],
            dir: './coverage',
            fixWebpackSourcePaths: true
        },
        port: 9876,
        colors: true,
        logLevel: config.LOG_INFO,
        autoWatch: true,
        browsers: ['PhantomJS'],
        singleRun: true,
        captureTimeout: 60000,
        browserDisconnectTimeout: 10000,
        browserNoActivityTimeout: 60000,
        browserDisconnectTolerance: 3
    });
};

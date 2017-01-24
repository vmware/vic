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

module.exports = {
    devtool: 'inline-source-map',
    verbose: true,
    resolve: {
        extensions: ['', '.ts', '.js', '.scss', '.html'],
        modulesDirectories: ['node_modules', 'src']
    },
    module: {
        loaders: [
            {
                test: /\.ts$/,
                loaders: ['awesome-typescript-loader?sourceMap=false&inlineSourceMap=true', 'angular2-template-loader']
            },
            {test: /\.json$/, loader: 'json-loader'},
            {test: /\.html$/, loader: 'html-loader'},
            {test: /\.scss$/, loaders: ['to-string-loader', 'css-loader', 'sass-loader']},
            {test: /\.css$/, loaders: ['to-string-loader', 'css-loader']}
        ],
        postLoaders: [
            {
                test: /\.ts$/,
                loader: 'istanbul-instrumenter-loader',
                exclude: [
                    'node_modules',
                    /karma-entry\.ts$/,
                    /\.(e2e|spec)\.ts$/
                ]
            }
        ]
    }
}
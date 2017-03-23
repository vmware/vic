var path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const ChunkManifestPlugin = require('chunk-manifest-webpack-plugin');
const WebpackChunkHash = require('webpack-chunk-hash');
const webpack = require('webpack');

function getRoot() {
    return path.resolve(__dirname, '../');
}

function getPublicPath(env) {
    if (!env) {
        return null;
    }
    return env.is_production ?
        '/ui/vic/resources/dist/' :
        env.is_standalone ?
        '/' :
        '/ui/vic/resources/build-dev/';
}

function getPath(env) {
    if (!env) {
        return null;
    }
    return env.is_production ?
        path.resolve(getRoot(), '../main/webapp/resources/dist/') :
        path.resolve(getRoot(), '../main/webapp/resources/build-dev/');
}

module.exports = function(env) {
    const base = {
        context: path.resolve(__dirname, '../'),
        devtool: env.is_production ? 'source-map' : 'cheap-module-source-map',
        entry: {
            'vendor': './src/vendor.ts',
            'main': env.is_aot ? './src/main-aot.ts' : './src/main.ts'
        },
        output: {
            path: getPath(env),
            filename: '[name].[chunkhash].js',
            chunkFilename: '[name].[chunkhash].js',
            publicPath: getPublicPath(env)
        },
        resolve: {
            modules: ['node_modules', 'src'],
            extensions: [".ts", ".js", ".scss", ".css", ".json", ".html"]
        },
        module: {
            rules: [
                {
                    test: /\.js$/,
                    enforce: 'pre',
                    loader: 'source-map-loader',
                    exclude: [
                        path.resolve(__dirname, 'node_modules/rxjs'),
                        path.resolve(__dirname, 'node_modules/@angular'),
                        /node_modules\/mutationobserver-shim\/dist/,
                        /node_modules\/clarity-icons/,
                        /node_modules\/clarity-angular/
                    ]
                },
                {
                    test: /\.ts$/,
                    enforce: 'pre',
                    loader: 'tslint-loader',
                    exclude: [
                        /node_modules\/clarity-ui/,
                        /node_modules\/clarity-icons/,
                        /node_modules\/clarity-angular/,
                        /aot/
                    ]
                },
                {
                    test: /\.ts$/,
                    use: [
                        'awesome-typescript-loader',
                        'angular2-template-loader'
                    ]
                },
                {
                    test: /\.css$/,
                    use: [
                        'to-string-loader',
                        'css-loader'
                    ]
                },
                {
                    test: /\.html$/,
                    loader: 'raw-loader',
                    exclude: [
                        path.resolve(__dirname, 'src/index.html')
                    ]
                }
            ]
        },
        plugins: [
            new webpack.ContextReplacementPlugin(
                /angular(\\|\/)core(\\|\/)(esm(\\|\/)src|src)(\\|\/)linker/,
                __dirname
            ),
            new HtmlWebpackPlugin({
                template: 'src/index.html',
                chunksSortMode: 'dependency'
            }),
            new ExtractTextPlugin({
                filename: 'clarity-ui-min.css',
                disable: false,
                allChunks: true
            }),
            new webpack.optimize.CommonsChunkPlugin({
                name: ['vendor', 'manifest'],
                minChunks: Infinity
            }),
            new CopyWebpackPlugin([
                {
                    from: path.resolve(getRoot(), 'src/assets'),
                    to: path.resolve(getPath(env), 'assets')
                },{
                    from: path.resolve(getRoot(), 'node_modules/clarity-ui/clarity-ui.min.css'),
                    to: path.resolve(getPath(env))
                }], {
                    ignore: ['*.scss', 'index.html']
                }
            )
        ],
        performance: {
            hints: false
        },
        target: 'web',
        stats: 'errors-only'
    };

    base.module.rules.splice(0, 0, {
        test: /environment\.ts$/,
        enforce: 'pre',
        loader: 'string-replace-loader',
        options: {
            search: /production\:.*/,
            replace: env.is_production ? 'production: true' : 'production: false'
        }
    },{
        test: /index\.html$/,
        enforce: 'post',
        loader: 'string-replace-loader',
        options: {
            search: 'APP_BASE_URL',
            replace: getPublicPath(env)
        }
    });

    if (!env.is_production) {
        base.devServer = {
            port: 3001,
            host: '0.0.0.0',
            historyApiFallback: true,
            watchOptions: {
                aggregateTimeout: 300,
                poll: 1000
            },
            outputPath: path.resolve(__dirname, '../../main/webapp/resources/build-dev')
        };
        base.module.rules.push({
            test: /\.scss$/,
            use: [
                'raw-loader',
                'sass-loader'
            ],
            exclude: [/node_modules/]
        });
    } else {
        base.module.rules.splice(0, 0, {
            test: /\.ts$/,
            loader: 'string-replace-loader',
            options: {
                search: /console\.log.*/,
                replace: '',
                flags: 'g'
            },
            include: [/src/]
        },{
            test: /\.scss$/,
            use: [
                'style-loader',
                'css-loader',
                'sass-loader'
            ],
            exclude: [/node_modules/]
        });
        base.plugins.splice(1, 0,
            new webpack.optimize.UglifyJsPlugin()
        );
        base.plugins.splice(5, 0,
            new webpack.HashedModuleIdsPlugin(),
            new WebpackChunkHash(),
            new ChunkManifestPlugin({
                filename: 'chunk-manifest.json',
                manifestVariable: 'webpackManifest'
            })
        );
    }

    return base;
};

// Include gulp
var gulp = require('gulp');

// Include Our Plugins
var jshint = require('gulp-jshint');
var sass = require('gulp-sass');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var uglifycss = require('gulp-uglifycss');
var rename = require('gulp-rename');
var inject = require('gulp-inject');
var bowerFiles = require('main-bower-files');
var angularFilesort = require('gulp-angular-filesort');
var server = require('gulp-server-livereload');
var clean = require('gulp-clean');
var runSequence = require('run-sequence');
var gulpNgConfig = require('gulp-ng-config');
var del = require('del');

//Configuration
var finalBuildDestination = '../server/client/'
var buildDestination = 'build';
var devDestination = 'dev';
var appFolder = 'app';

/******
 * DEV
*******/

gulp.task('clean-dev', function () {
	return gulp.src(devDestination, {read: false})
		.pipe(clean());
});

//Inject vendor
gulp.task('inject', function () {
  return gulp.src(appFolder+'/index.html', {read:false})
        .pipe(inject(gulp.src(bowerFiles()), {name: 'bower'}))
        .pipe(inject(
            gulp.src(appFolder+'/**/*.js') // gulp-angular-filesort depends on file contents, so don't use {read: false} here
            .pipe(angularFilesort())
        ))
        .pipe(gulp.dest(devDestination));
});

// Watch Files For Changes
gulp.task('watch', function() {
  gulp.watch(appFolder+'/**/*.js',['inject']);
});

gulp.task('webserver', function() {
  gulp.src(devDestination)
    .pipe(server({
      livereload: true,
      directoryListing: false,
      open: true
    }));
});

gulp.task('build-server', function() {
  gulp.src(buildDestination)
    .pipe(server({
      livereload: true,
      directoryListing: false,
      open: true
    }));
});

gulp.task('compose-server', function() {
  gulp.src(buildDestination)
    .pipe(server({
      livereload: true,
      directoryListing: false,
      open: true,
      host: "0.0.0.0"
    }));
});

gulp.task('delete-constant', function () {
  return del([appFolder+'/scripts/constants.js'])
});

//Creating constants for dev environment
gulp.task('dev-constants', ['delete-constant'], function () {
  return gulp.src('config.json')
    .pipe(gulpNgConfig('TasksApp', {
      environment: 'local',
      createModule: false
    }))
    .pipe(rename('constants.js'))
    .pipe(gulp.dest(appFolder+'/scripts/',{force: true}))
});

gulp.task('compose-constants', ['delete-constant'], function () {
  return gulp.src('config.json')
    .pipe(gulpNgConfig('TasksApp', {
      environment: 'compose',
      createModule: false
    }))
    .pipe(rename('constants.js'))
    .pipe(gulp.dest(appFolder+'/scripts/',{force: true}))
});

gulp.task('build-watch', function() {
  gulp.watch(appFolder+'/**/*.js',['build-js']);
  gulp.watch(appFolder+'/**/*.html',['buildinject']);
  gulp.watch(appFolder+'/**/*.css',['build-css']);
});

/******
 * BUILD
*******/

gulp.task('clean-build', function () {
	return gulp.src(buildDestination, {read: false})
		.pipe(clean());
});

//Creating constants for dev environment
gulp.task('deploy-constants',  ['delete-constant'], function () {
  return gulp.src('config.json')
    .pipe(gulpNgConfig('TasksApp', {
      environment: 'production',
      createModule: false
    }))
    .pipe(rename('constants.js'))
    .pipe(gulp.dest(appFolder+'/scripts/',{force: true}))
});

//Concatenates bower plugins js files and minimize directly
gulp.task('vendor-build-js', function () {
  var javascript = [];
  var regex = /\.js$/;
  bowerFiles().forEach(function (file) {
    if (regex.test(file)) {
      javascript.push(file);
    }
  });
  return gulp.src(javascript)
    .pipe(concat('vendor.min.js'))
    .pipe(uglify())
    .pipe(gulp.dest(buildDestination+'/vendor'));
});

//Concatenates bower plugins css files and minimize directly
gulp.task('vendor-build-css', function () {
  var css = [];
  var regex = /\.css$/;
  bowerFiles().forEach(function (file) {
    if (regex.test(file)) {
      css.push(file);
    }
  });
  return gulp.src(css)
    .pipe(concat('vendor.min.css'))
    .pipe(uglifycss({
      "maxLineLen": 80,
      "uglyComments": true
    }))
    .pipe(gulp.dest(buildDestination+'/vendor'));
});

//Concatenates js files and minimize directly
gulp.task('build-js', function () {
  return gulp.src(appFolder+'/**/*.js')
    .pipe(angularFilesort())
    .pipe(concat('app.min.js'))
    .pipe(uglify())
    .pipe(gulp.dest(buildDestination+'/js/'));
});

//Concatenates js files and minimize directly
gulp.task('build-css', function () {
  return gulp.src(appFolder+'/**/*.css')
    .pipe(concat('app.min.css'))
    .pipe(uglifycss())
    .pipe(gulp.dest(buildDestination+'/css/'));
});

//Copy media
gulp.task('build-media', function () {
  return gulp.src(appFolder+'/media/*')
    .pipe(gulp.dest(buildDestination+'/media/'));
});

//Copies index.html to destination folder
gulp.task('build-copy', function () {
  return gulp.src([appFolder+'/index.html', appFolder+'/**/*.html'])
    // .pipe(gulp.dest(buildDestination))
    // .pipe(gulp.src(appFolder+'/**/*.html'))
    .pipe(gulp.dest(buildDestination));
});

//Inject stuff into the index
gulp.task('buildinject',['build-copy'], function () {
  return gulp.src(buildDestination+'/index.html')
    .pipe(inject(gulp.src([buildDestination+'/vendor/*.js',buildDestination+'/vendor/*.css']), {name: 'bower', relative:true}))
    .pipe(inject(gulp.src([buildDestination+'/js/*.js',buildDestination+'/css/*.css']),{relative:true}))
    .pipe(gulp.dest(buildDestination));
});

gulp.task('clean-final-build', function () {
	return gulp.src(finalBuildDestination+'/*', {read: false})
		.pipe(clean({force:true}));
});

gulp.task('copy-build-final', function () {
  return gulp.src(buildDestination+'/**/*',{base: buildDestination})
    .pipe(gulp.dest(finalBuildDestination));
});


/***** IMPORTANT TASKS *****/

//Default task - Launch webserver
gulp.task('default', function (done) {
  runSequence('dev-constants','rebuild', 'build-server', 'build-watch',done);
});

//Default task - Launch webserver
gulp.task('compose', function (done) {
  runSequence('compose-constants','rebuild', 'compose-server', 'build-watch',done);
});

//Rebuild a package version ready to deploy without removing the previous files
gulp.task('rebuild', function (done) {
  runSequence(['vendor-build-js', 'vendor-build-css', 'build-js', 'build-css', 'build-media'], 'buildinject', done);
})

//Builds a package version ready to deploy
gulp.task('build', function(done) {
  runSequence('clean-build','deploy-constants','rebuild', done);
});

//Builds and deploys inside the server/client folder
gulp.task('deploy',['build'], function (done){
  runSequence('clean-final-build', 'copy-build-final', done);
});

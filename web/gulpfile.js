var gulp       = require('gulp');
// var beep       = require('beepbeep')
var gutil      = require('gulp-util');
var plumber    = require('gulp-plumber');
var uglify     = require('gulp-uglifyjs');
var sass       = require('gulp-ruby-sass');
var livereload = require('gulp-livereload');
var connect    = require('gulp-connect');
var htmlMin    = require('gulp-htmlmin')

var onError = function (err) {
    // beep([0, 0, 0]);
    gutil.log(gutil.colors.green(err));
};


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////
//////////     WEBSITE TASKS
//////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


gulp.task('connect', function() {
  connect.server({
    root: 'build',
    port: 8000,
    livereload: true
  });
});

// JS
gulp.task('uglifyjs', function() {
    return gulp.src([
        './bower_components/jquery/dist/jquery.min.js',
        './js/main.js'
    ])
    .pipe(plumber({
        errorHandler: onError
    }))
    .pipe(uglify('app.js', {
        compress: false
    }))
    .pipe(gulp.dest('./js/'))
    // .pipe(livereload());
    .pipe(connect.reload());
});

// Sass
gulp.task('sass', function() {
    return gulp.src([
        './scss/app.scss'
    ])
    .pipe(plumber({
        errorHandler: onError
    }))
    .pipe(sass({
        style: 'compressed',
        cacheLocation: './cache/.sass-cache'
    }))
    .pipe(gulp.dest('./css/'))
    // .pipe(livereload());
    .pipe(connect.reload());
});

// HTML
gulp.task('html', function() {
    return gulp.src([
        './src/*.html'
    ])
    .pipe(htmlMin({collapseWhitespace: true}))
    .pipe(gulp.dest('build/'))
    // .pipe(livereload());
    .pipe(connect.reload());
});


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////
//////////     WATCH AND BUILD TASKS
//////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Primary task to watch other tasks
// gulp.task('yo', function() {
gulp.task('watch', function () {
    // LiveReload
    // livereload.listen();

    // Watch JS
    gulp.watch('./js/main.js', ['uglifyjs']);

    // Watch Sass
    gulp.watch(['./scss/_mixins.scss', './scss/_styles.scss', './scss/app.scss'], ['sass']);

    // Watch HTML and livereload
    gulp.watch('./src/*.html', ['html']);
});

gulp.task('default', ['connect', 'watch']);

// gulp.task('watch', function () {
//   gulp.watch(['./app/*.html'], ['html']);
// });

// Manually build all
gulp.task('build', function() {
    gulp.start('uglifyjs', 'sass');
});

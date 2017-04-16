// See: https://markgoodyear.com/2014/01/getting-started-with-gulp/
var gulp = require('gulp'),
    concat = require('gulp-concat'),
    browserSync = require('browser-sync').create(),
    uglify = require('gulp-uglify'),
    cleanCSS = require('gulp-clean-css'),
    autoprefixer = require('gulp-autoprefixer'),
    htmlMin = require('gulp-htmlmin')

gulp.task('browserSync', function() {
   browserSync.init({
      server: {
         baseDir: 'build'
      },
   })
})

gulp.task('html', function() {
    gulp.src(['src/**.html'])
    .pipe(htmlMin({collapseWhitespace: true}))
    .pipe(gulp.dest('build/'))
    .pipe(browserSync.reload({
      stream: true
   }))
});

gulp.task('styles', function() {   
   gulp.src(['src/styles/*.css'])
   .pipe(concat('style.css'))
   .pipe(autoprefixer('last 2 versions'))
   .pipe(cleanCSS())
   .pipe(gulp.dest('build/styles/'))
   .pipe(browserSync.reload({
      stream: true
   }))
});

gulp.task('scripts', function() {
  return gulp.src('src/scripts/**/*.js')
    // .pipe(jshint('.jshintrc'))
    // .pipe(jshint.reporter('default'))
    .pipe(concat('main.js'))
    // .pipe(gulp.dest('build/scripts/'))
    // .pipe(rename({suffix: '.min'}))
    .pipe(uglify())
    .pipe(gulp.dest('build/scripts/'))
    // .pipe(notify({ message: 'Scripts task complete' }));
});

gulp.task('default', ['browserSync', 'styles', 'html', 'scripts'], function (){
   gulp.watch('src/styles/*.css', ['styles']);
   gulp.watch('src/**.html', ['html']);
   gulp.watch('src/scripts/**/*.js', ['scripts']);
});   
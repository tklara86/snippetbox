
var gulp = require("gulp"),
    sass = require("gulp-sass"),
    postcss = require("gulp-postcss"),
    autoprefixer = require("autoprefixer"),
    cssnano = require("cssnano"),
    sourcemaps = require("gulp-sourcemaps");




var paths = {
    styles: {
        src: "src/scss/**/*.scss",
        dest: "../css/"
    }


};

function style() {
    return gulp
        .src(paths.styles.src)
        .pipe(sourcemaps.init())
        .pipe(sass())
        .on("error", sass.logError)
        .pipe(postcss([autoprefixer(), cssnano()]))
        .pipe(sourcemaps.write('maps'))
        .pipe(gulp.dest(paths.styles.dest))
}

function watch() {

    gulp.watch(paths.styles.src, style);

}


exports.watch = watch
exports.style = style;


var build = gulp.parallel(style, watch);
gulp.task('default', build);


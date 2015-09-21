$(document).ready(function() {
    $('pre code').each(function(i, block) {
        $(block).addClass("line-numbers");
        Prism.highlightElement(block, false)
    });
});
$(document).ready(function () {
    $('pre code').each(function (i, block) {
        $(block).addClass("line-numbers");
        Prism.highlightElement(block, false)
    });
    $(".scroll-to").click(function (event) {
        event.preventDefault();
        console.log("scroll to:", $(event.target).attr("href"));
        //get the top offset of the target anchor
        var target_offset = $($(event.target).attr("href")).offset();
        var target_top = target_offset.top;
        $('html, body').animate({scrollTop: target_top}, 500, 'easeInSine');
    });
});

/*
 * Code to check that usernames are unique.
 *
 * Largely based on: http://grinninggecko.com/easy-ajax-jquery/ by Garth Gutenberg.
 */


/*
 * noConflict() required in case the existing site uses another library that includes $()
 */
$j = jQuery.noConflict();
 
/*
 * Standard jQuery (document).ready()
 * Begins processing the enclosed function once the DOM is loaded
 */
$j(document).ready(function() {
    function checkUid(uid) {
        $j.post(
            'checkuid?uid=' + uid,
            function(data) {
                obj = JSON.parse(data);
                alert('Username ' + obj.Uid + ' is ' + obj.Available);
            }
        );
    }
 
    $j('#uid').bind('blur', function() {
        var username = $j(this).val();
        checkUid(username);
    });
});
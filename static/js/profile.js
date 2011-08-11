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
				if (uid == "") {
					$j('#uid_required').show();
					$j('#uid_not_ok').hide();
					$j('#uid_ok').hide();
					return;
				}
				else if (obj.Available == "available") {
					$j('#uid_required').hide();
					$j('#uid_not_ok').hide();
					$j('#uid_ok').show();
				} else {
					$j('#uid_required').hide();
					$j('#uid_ok').hide();
					$j('#uid_not_ok').show();					
				}
            }
        );
    }
 
    $j('#uid').bind('blur', function() {
        var username = $j(this).val();
        checkUid(username);
    });
});
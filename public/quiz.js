var debugOn = false;
var save_button_text = "Save Progress";
var current_question = null;
var assignments = {};

var bool_only_dupes = false;
var bool_hide_answered = false;
var bool_all_answers = false;
var bool_saving = false;

var assign = function(question, answer, check) {
	// Question is a jquery object
	// Answer is a string "a_###"
	// Check is a bool (should we reevaluate dupes or not---only required
	// if this was called by the user, not by load_state()
	check = (typeof check === "undefined") ? true : check;

	id = answer.attr('id');

	var q_id = question.attr('id');
	var content = answer.html();

	// If the question has an existing assignment, remove it.
	// Also, check if the change affects the status of any duplicates.
	// reevaluate_dupes checks that.
	// Will only be executed if this function is called by clicking on an answer line
	if (check && assignments[q_id]) {
		var a_id = assignments[q_id];
		$('#'+q_id).removeClass("duplicate");
		if (bool_only_dupes) $('#'+q_id).addClass('vanish');
		delete assignments[q_id];
		reevaluate_dupes(a_id);
	}

	// scan_for_dupes function checks the assignments object
	// for existing entries that have been set to the same title id
	if (scan_for_dupes(id)) {
		answer.addClass("duplicate");
		question.addClass("duplicate");
	}

	answer.addClass("chosen");
	question.addClass("answered");
	if (bool_hide_answered) question.addClass('vanish');
	question.children('.q_title').html(content);
	//question.children('.q_title').addClass("answered");
	question.children('input').val(id.substring(2));

	assignments[question.attr('id')] = id;
	question.removeClass('current');
	$('.palette').removeClass('current');
	//question = null;
};

// Scan the assignments object for an existing entry
// which has the same value that's passed to the function.
var scan_for_dupes = function(id) {
	for (var picture in assignments) {
		if (assignments[picture] === id) {
			// Mark this div as class 'duplicate'
			$('#'+picture).addClass('duplicate');
			return true;
		}
	}
	return false;
};

// We only call this function if an answer id was duplicated, and a change
// might have removed the duplication.
var reevaluate_dupes = function(a_id) {
	var first_pic_id = false;
	for (var picture in assignments) {
		if (assignments[picture] === a_id) {
			if (first_pic_id !== false) return;
			// If we didn't return, then this is the first time we found the answer
			// id in question.
			first_pic_id = picture;
		}
	}
	// If we got all this way without returning, that means that there is exactly
	// one entry in assignments that contains the answer id. It's no longer a duplicate,
	// so we can remove the 'duplicate' classes.
	if (first_pic_id) {
		$('#'+first_pic_id).removeClass('duplicate');
		$('#'+a_id).removeClass('duplicate');
		if (bool_only_dupes) $('#'+first_pic_id).addClass('vanish');
	} else {
		// At this point, if nothing was found, that means at the answer id
		// is no longer used at all.
		$('#'+a_id).removeClass('chosen')
	}
};

var b64_alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_";
// Encode and Decode are way too wasteful with the size of their b64 strings.
// Using huge base-184 and base-64 numbers (see base.rb), a 183 element array
// can be encoded in 230 chars. JavaScript doesn't support numbers that big, though.
// Using 2 base-64 digits per base-184 digit, a 183 element array requires 366 chars.
// Poor form for something going into a URL.
var encode = function (status_array, qty) {
	var src_base = qty + 1;
	var dst_base = 64;
	var b64_string = ""
	for (var i = 0; i < status_array.length; i++) {
		c = status_array[i];
		b64_string += b64_alphabet[Math.floor(c/dst_base)];
		b64_string += b64_alphabet[c%dst_base];
	}
	return b64_string;
}

var decode = function (b64_string) {
	var dst_base = 64;
	var rebuilt = [];
	var i = 0;
	while (i < b64_string.length - 1) {
		value = b64_alphabet.indexOf(b64_string[i++]) * dst_base;
		value += b64_alphabet.indexOf(b64_string[i++]);
		if (debugOn) console.log("Decoded saved answer: " + value);
		rebuilt.push(value.toString());
	}
	return rebuilt;
};

var save_state = function() {
	// Yeah there is already an 'assignments' associative array, but
	// this array will inherently be sorted first to last, and have
	// entries for non-answered questions.
	$all_questions = $('.question');
	var q_status = [];
	var qty = $all_questions.length;

	$all_questions.each(function() {
		var value = $(this).children('input').attr('value');
		q_status.push(value);
	});

	var encoded = encode(q_status, qty);

	var rxp = /^[^?]+/; // For getting the base url before the ?param
	var base_url = rxp.exec(window.location.href)[0];
	base_url += "?" + encoded;
	//var fake_url = "http://www.spamusers.com/test.html?" + encoded;

	// Get ourselves a goo.gl address.
	var response = $.ajax({
		url: "https://www.googleapis.com/urlshortener/v1/url",
		type: "POST",
		data: JSON.stringify({ longUrl: base_url }),
		contentType: "application/json; charset=utf-8",
		dataType: "json",
		success: function(data) {
			if (data.id != null) {
				var shortened = data.id;
				console.log("Quiz state saved to " + shortened);
				var url_bar = $('.title.url');
				$('.title.text').addClass('vanish');
				$('#save').html('Saved! <img class="emote" src="graphics/icon_biggrin.png">');
				url_bar.html(shortened);
				url_bar.addClass('yet_still_visible');
			}
		},
		error: function(data) {
			console.log("Google hates you! Ya scum!");
			var url_bar = $('.title.url');
			$('.title.text').addClass('vanish');
			$('#save').html('Um, oops... <img class="emote" src="graphics/icon_mad.png">');
			url_bar.html("Google isn't in the mood right now.");
			url_bar.addClass('yet_still_visible');
		}
	});

};

var load_state = function(b64_string) {
	var array_status = decode(b64_string);

	// Now actually do the work (ugh)
	var i = 0;
	$('.question').each(function() {
		// Return early if we've run out of saved answers before we hit
		// the end of the .question query
		if (i == array_status.length) {
			console.log("Not enough saved answers to populate all the questions!");
			return;
		}
		if (array_status[i] != 0) {
			var a_id = "#a_"+array_status[i];
			var answer = $(a_id);
			// This next 'if' ensures that the decoded answer actually exists.
			// It always should unless someone types random ?params into the url,
			// like "index.html?blahblahblah"
			if (answer.length != 0)
				// The 'false' means that the function won't bother reevaluating
				// the assignments hash to see if a changed assignment removed a
				// duplicate.
				assign($(this), answer, false);
			else
				console.log("Ignoring invalid decoded answer " + a_id);
			i++;
		} else {
			if (debugOn) console.log("Ignoring unanswered question " + i);
			i++;
		}
	});

	if (i < array_status.length) {
		console.log("Too many answered questions. Some were ignored.");
	}
	//return array_status;
};

$(document).ready(function() {

	// If the page URL has a ?parameter, then automatically call load_state
	// on it.
	var load_rxp = /\?([\w\-]+)/; // For getting the ?param
	var load_rxp2 = /^.*?(?=\?)/; // For getting everything before the ?param
	if (load_rxp.test(window.location.href)) {
		var b64_string = load_rxp.exec(window.location.href)[1];
		console.log("Auto-loading saved state.");
		load_state(b64_string);
		//window.location.href = load_rxp2.exec(window.location.href)[0];
		window.history.replaceState({}, null, "?");
	}

	$('.question').click(function() {
		if (current_question) {
			// If there's already a current question, remove the class from it
			current_question.removeClass('current');
			// If current question is self, then we just want to deselect it
			// and leave current_question == null
			if ($(this).attr('id') === current_question.attr('id')) {
				current_question = null;
				// also 'deselect' the palette
				$('.palette').removeClass('current');
				return;
			}
		}
		current_question = $(this);
		current_question.addClass('current');
		$('.palette').addClass('current');
		if (debugOn) console.log("Clicked on picture id " + current_question.attr('id'));
	});

	$('.question.changeup').hover(function() {
		//Mouseover in
		target_img = $(this).children('img')[0];
		target_img.src = target_img.src.replace(/\.(?=\w{3}$)/,"a.");
	}, function () {
		//Mouseover out
		target_img = $(this).children('img')[0];
		target_img.src = target_img.src.replace(/a\.(?=\w{3}$)/,".");
	});

	$('.palette #CLEAR').click(function() {
		if (current_question !== null) {
			var q_id = current_question.attr('id');
			if (assignments[q_id]) {
				if (debugOn) console.log("Clearing assignment of "+ q_id);
				var a_id = assignments[q_id];
				quest = $('#'+q_id);
				quest.removeClass("duplicate answered current");

				// This change may also change this question's visibility...
				if (bool_hide_answered) quest.removeClass('vanish');
				else if (bool_only_dupes) quest.addClass('vanish');

				quest.children('.q_title').removeClass('answered');
				quest.children('.q_title').html("NO TITLE SELECTED");
				quest.children('input').val("0");
				$('.palette').removeClass('current');
				delete assignments[q_id];
				reevaluate_dupes(a_id);
			}
		}
	});

	$('.palette a:not(#CLEAR)').click(function() {
		if (current_question !== null) {
			if (debugOn) console.log("Clicked on title id   " + id);
			assign(current_question, $(this));
			current_question = null;
		}
	});

	$('#a_show_all').click(function() {
		$(this).toggleClass('depressed');
		$('.palette a').toggleClass('yet_still_visible');
	});

	$('#q_toggle_answered').click(function() {
		if (bool_only_dupes) {
			$('#q_only_dupes').removeClass('depressed');
			$('.question.vanish').removeClass('vanish');
			bool_only_dupes = false;
		}
		if (bool_hide_answered) {
			$(this).removeClass('depressed');
			bool_hide_answered = false;
			$('.question.answered').removeClass('vanish');
		} else {
			$(this).addClass('depressed');
			bool_hide_answered = true;
			$('.question.answered').addClass('vanish');
		}
	});

	$('#q_only_dupes').click(function() {
		//First disable everything that the toggle_answered button does, in
		//case it's been turned on.
		if (bool_hide_answered) {
			$('#q_toggle_answered').removeClass('depressed');
			$('.question.answered').removeClass('vanish');
			bool_hide_answered = false;
		}
		if (bool_only_dupes) {
			$(this).removeClass('depressed');
			bool_only_dupes = false;
			$('.question').removeClass('vanish');
		} else {
			$(this).addClass('depressed');
			bool_only_dupes = true;
			$('.question:not(.duplicate)').addClass('vanish');
		}
	});

	$('#save').click(function() {
		if (bool_saving) {
			bool_saving = false;
			$(this).removeClass('depressed');
			$(this).html(save_button_text);
			$('.title.url').removeClass('yet_still_visible');
			$('.title.text').removeClass('vanish');
		} else {
			bool_saving = true;
			save_button_text = $(this).html();
			$(this).html("Saving . . .");
			$(this).addClass('depressed');
			save_state();
		}
	});

});

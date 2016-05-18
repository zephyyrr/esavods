ESAVods = function($) {

	var fillEvents = function() {
		var events = [
			{name: "ESA2016", dates: {start: "2016-06-23T20:04:00Z", end: "2016-06-29T20:04:00Z"} },
			{name: "ESA2015", dates: {start: "2016-07-26T14:00:00Z", end: "2016-08-07T15:03:00Z"} },
			{name: "ESA2014", dates: {start: "2016-06-24T16:00:00Z", end: "2016-07-02T13:37:00Z"} }
		];
		events.forEach(function(event, i) {
			$.ajax({
				type: "POST",
				url: "/events",
				data: JSON.stringify(events[i]),
				contentType: "application/json",
				dataType: "json"
			});
		});
	}

	var fillRuns = function() {
		var runs = [
			{game: "Golden Sun", players:["BaalNocturno"], category: "100%", type: "normal", console: "Gameboy Advance", length: (2*3600+10*60+4)*1000, event: "ESA2014", tags: ["test", "example"], vods: {"youtube": "https://www.youtube.com/watch?v=jvhrfc4_PBQ"}},
			{game: "Final Fantas: Crystal Chronicles", players:["Cereth", "Grokken", "MLSTRM", "Neviutz"], category: "All Dungeons", type: "co-op", console: "Gamecube", length: (1*3600+40*60+0)*1000, event: "ESA2015", tags: ["test", "demo"], vods: {"youtube": "https://www.youtube.com/watch?v=wJMdyEftGJU"}},
			{game: "Castle of Illusion HD", players:["Grukk", "KrazyRasmus", "MaxieTheHatter","Pogington"], category: "any%", type: "race", console: "PC", length: (0*3600+32*60+5)*1000, event: "ESA2015", tags: ["test", "example"], vods: {"youtube": "https://www.youtube.com/watch?v=jQZQrWaLMO0"}}
		];
		runs.forEach(function(events, i) {
			$.ajax({
				type: "POST",
				url: "/runs",
				data: JSON.stringify(runs[i]),
				contentType: "application/json",
				dataType: "json"
			});
		});
	}

	function fillDB() {
		fillEvents();
		fillRuns();
	}

	return {
		fillDB: fillDB
	}
}($)

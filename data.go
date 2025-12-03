package main

var config = Config{
	Title:   "][ nobloat.org",
	Slogan:  "pragmatic software minimalism",
	BaseURL: "https://nobloat.org",
	Links: map[string]string{
		"Choosing boring technology":                         "https://boringtechnology.club/",
		"Radical simplicity":                                 "https://www.radicalsimpli.city/",
		"Frameworkless movement":                             "https://frameworklessmovement.org/",
		"No frameworks Part II - Matteo Vaccari":             "http://matteo.vaccari.name/blog/archives/1022",
		"Local first software":                               "https://www.inkandswitch.com/local-first/",
		"Local-first software: Talk by Peter Van Hardenberg": "https://www.youtube.com/watch?v=KrPsyr8Ig6M",
		"Minimal viable program":                             "https://joearms.github.io/published/2014-06-25-minimal-viable-program.html",
		"YAGNI (You ain't gonna need it)":                    "https://en.wikipedia.org/wiki/You_ain%27t_gonna_need_it",
		"suckless.org":                                       "https://suckless.org/philosophy/",
		"cat -v considered harmful":                          "https://harmful.cat-v.org/",
		"Simplifier":                                         "https://simplifier.neocities.org/",
		"zeitkapsl.eu":                                       "https://zeitkapsl.eu",
		"hardcode.at":                                        "https://hardcode.at",
		"spiessknafl.at/peter":                               "https://spiessknafl.at/peter",
	},
	Projects: map[string]string{
		"[nobloat/css](https://github.com/nobloat/css)":                                           "modular vanilla CSS3 components",
		"[nobloat/bare-jvm](https://github.com/nobloat/bare-jvm)":                                 "[baremessages](https://baremessages.org/) implementation for the JVM",
		"[cinemast/dbolve](https://github.com/cinemast/dbolve)":                                   "database migration tool for golang projects",
		"[nobloat/tinyviper](https://github.com/nobloat/tinyviper)":                               "alternative to the famous [spf13/viper](https://github.com/spf13/viper) configuration library",
		"[oliverselinger/failsafe-executor](https://github.com/oliverselinger/failsafe-executor)": "distributed task scheduler for Java 11+ projects (lightweight alternative to quartz)",
		"[oliverselinger/db-evolve](https://github.com/oliverselinger/db-evolve)":                 "database migration tool for Java 11+ projects (lightweight alternative to liquibase or flyway)",
		"[nobloat/svelte-router](https://github.com/nobloat/svelte-router)":                       "router for svelte, in case one does not want to use SvelteKit",
	},
	Tools: []Tool{
		{Name: "bundlephobia", Description: "A tool to analyze the size of your JavaScript packages", URL: "https://bundlephobia.com/"},
	},
}

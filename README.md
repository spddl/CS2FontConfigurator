# CS2FontConfigurator
[![Downloads][1]][2] [![GitHub stars][3]][4]

[1]: https://img.shields.io/github/downloads/spddl/CS2FontConfigurator/total.svg
[2]: https://github.com/spddl/CS2FontConfigurator/releases "Downloads"

[3]: https://img.shields.io/github/stars/spddl/CS2FontConfigurator.svg
[4]: https://github.com/spddl/CS2FontConfigurator/stargazers "GitHub stars"

### [Download](https://github.com/spddl/CS2FontConfigurator/releases)

Like CSGO, CS2 uses a [FontConfig](https://www.freedesktop.org/wiki/Software/fontconfig/), the difference is that Valve provides a preset file for CS2 that shows how fonts should be replaced.
> core\panorama\fonts\conf.d\42-repl-global.conf

We use this advantage.

We only change this file and the first loaded `fonts.conf`, otherwise fonts installed in the user area under Windows will not be found.
> csgo\panorama\fonts\fonts.conf

![screenshot](https://i.imgur.com/rN2zzOh.png)
# CS2FontConfigurator

[![Downloads][1]][2] [![GitHub stars][3]][4]

[1]: https://img.shields.io/github/downloads/spddl/CS2FontConfigurator/total.svg
[2]: https://github.com/spddl/CS2FontConfigurator/releases "Downloads"
[3]: https://img.shields.io/github/stars/spddl/CS2FontConfigurator.svg
[4]: https://github.com/spddl/CS2FontConfigurator/stargazers "GitHub stars"

### [Download](https://github.com/spddl/CS2FontConfigurator/releases)

Like CSGO, CS2 uses [FontConfig](https://www.freedesktop.org/wiki/Software/fontconfig/) which is not widely used on Windows at all.
It seems that it is not fully supported or maybe even modified. Since the configuration is simply wrong according to the [official documentation of Fontconfig](https://fontconfig.pages.freedesktop.org/fontconfig/fontconfig-user.html)

```conf
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 39: unknown element "fontpattern"
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 40: unknown element "fontpattern"
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 41: unknown element "fontpattern"
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 42: unknown element "fontpattern"
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 43: unknown element "fontpattern"
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\csgo\panorama\fonts\fonts.conf", line 86: saw string, expected range
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\core\panorama\fonts\../../../core/panorama/fonts/conf.d/41-repl-os-win.conf", line 148: Having multiple values in <test> isn't supported and may not work as expected
Fontconfig warning: "C:\Steam\steamapps\common\Counter-Strike Global Offensive\game\core\panorama\fonts\../../../core/panorama/fonts/conf.d/41-repl-os-win.conf", line 160: Having multiple values in <test> isn't supported and may not work as expected
```

It is hard to tell if this is planned or will change in the future.
My solution is to fix the `fontconfig.conf` and apply my target font with an additional entry to all texts.
No valve font will be deleted and a backup of the `fontconfig.conf` will be created first.

![screenshot](https://i.imgur.com/B2fZNWL.png)

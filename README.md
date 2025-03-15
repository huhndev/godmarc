<p align="center">
  <code>godmarc</code> is a TUI DMARC report analyzer.
</p>

#

### Main Features

- Parse and display DMARC reports
- Mutt-inspired interface
- View individual reports or an aggregated summary

#

> [!IMPORTANT]
> `godmarc` is pre-alpha software.

> [!NOTE]
> In [CHANGELOG.md] you can follow recent changes.
> [ROADMAP.md] shows our future plans.

> [!CAUTION]
> `godmarc` is part of a [vibe coding] project and was created with the help of [Claude 3.7 Sonnet].

***

### Usage

1. Place your DMARC reports in `~/.godmarc`
2. Start `godmarc`

### DMARC Report Format

`godmarc` expects DMARC reports in the standard XML format. These reports are typically received via email when you have configured DMARC for your domain.

Example snippet of a DMARC report:

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<feedback>
  <report_metadata>
    <org_name>google.com</org_name>
    <email>noreply-dmarc-support@google.com</email>
    <report_id>1234567890123456789</report_id>
    <date_range>
      <begin>1596240000</begin>
      <end>1596326399</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>none</p>
    <sp>none</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <!-- Record data -->
  </record>
</feedback>
```

### Contributing

See [CONTRIBUTING.md]

### Security

See [SECURITY.md]

### License

The package may be used under the terms of the ISC License a copy of
which may be found in the file [LICENSE].

Unless you explicitly state otherwise, any contribution submitted for inclusion
in the work by you shall be licensed as above, without any additional terms or
conditions.

[ROADMAP.md]: https://github.com/huhnsystems/godmarc/blob/master/docs/ROADMAP.md
[CHANGELOG.md]: https://github.com/huhnsystems/godmarc/blob/master/docs/CHANGELOG.md
[CONTRIBUTING.md]: https://github.com/huhnsystems/godmarc/blob/master/docs/CONTRIBUTING.md
[SECURITY.md]: https://github.com/huhnsystems/godmarc/blob/master/docs/SECURITY.md
[LICENSE]: https://github.com/huhnsystems/godmarc/blob/master/LICENSE
[vibe coding]: https://en.wikipedia.org/wiki/Vibe_coding
[Claude 3.7 Sonnet]: https://en.wikipedia.org/wiki/Claude_(language_model)#Claude_3.7

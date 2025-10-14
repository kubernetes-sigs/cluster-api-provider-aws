<!-- omit in toc -->
# goteamsnotify

A package to send messages to a Microsoft Teams channel.

[![Latest release][githubtag-image]][githubtag-url]
[![Go Reference][goref-image]][goref-url]
[![License][license-image]][license-url]
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/go-teams-notify)](https://github.com/atc0005/go-teams-notify)
[![Lint and Build](https://github.com/atc0005/go-teams-notify/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/go-teams-notify/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/go-teams-notify/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/go-teams-notify/actions/workflows/project-analysis.yml)

<!-- omit in toc -->
## Table of contents

- [Project home](#project-home)
- [Overview](#overview)
- [Features](#features)
- [Project Status](#project-status)
- [Supported Releases](#supported-releases)
- [Changelog](#changelog)
- [Usage](#usage)
  - [Add this project as a dependency](#add-this-project-as-a-dependency)
  - [Webhook URLs](#webhook-urls)
    - [Expected format](#expected-format)
    - [How to create a webhook URL (Connector)](#how-to-create-a-webhook-url-connector)
  - [Examples](#examples)
    - [Basic](#basic)
    - [Specify proxy server](#specify-proxy-server)
    - [User Mention](#user-mention)
    - [Tables](#tables)
    - [Set custom user agent](#set-custom-user-agent)
    - [Add an Action](#add-an-action)
    - [Toggle visibility](#toggle-visibility)
    - [Disable webhook URL prefix validation](#disable-webhook-url-prefix-validation)
    - [Enable custom patterns' validation](#enable-custom-patterns-validation)
- [Used by](#used-by)
- [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/go-teams-notify) for the
latest code, to file an issue or submit improvements for review and potential
inclusion into the project.

## Overview

The `goteamsnotify` package (aka, `go-teams-notify`) allows sending messages
to a Microsoft Teams channel. These messages can be composed of legacy
[`MessageCard`][msgcard-ref] or [`Adaptive Card`][adaptivecard-ref] card
formats.

Simple messages can be created by specifying only a title and a text body.
More complex messages may be composed of multiple sections (`MessageCard`) or
containers (`Adaptive Card`), key/value pairs (aka, `Facts`) and externally
hosted images. See the [Features](#features) list for more information.

**NOTE**: `Adaptive Card` support is currently limited. The goal is to expand
this support in future releases to include additional features supported by
Microsoft Teams.

## Features

- Submit simple or complex messages to Microsoft Teams
  - simple messages consist of only a title and a text body (one or more
    strings)
  - complex messages may consist of multiple sections (`MessageCard`),
    containers (`Adaptive Card`) key/value pairs (aka, `Facts`) and externally
    hosted images
- Support for Actions, allowing users to take quick actions within Microsoft
  Teams
  - [`MessageCard` `Actions`][msgcard-ref-actions]
  - [`Adaptive Card` `Actions`][adaptivecard-ref-actions]
- Support for [user mentions][adaptivecard-user-mentions] (`Adaptive
  Card` format)
- Configurable validation of webhook URLs
  - enabled by default, attempts to match most common known webhook URL
    patterns
  - option to disable validation entirely
  - option to use custom validation patterns
- Configurable validation of `MessageCard` type
  - default assertion that bare-minimum required fields are present
  - support for providing a custom validation function to override default
    validation behavior
- Configurable validation of `Adaptive Card` type
  - default assertion that bare-minimum required fields are present
  - support for providing a custom validation function to override default
    validation behavior
- Configurable timeouts
- Configurable retry support

## Project Status

In short:

- The upstream project is no longer being actively developed or maintained.
- This fork is now a standalone project, accepting contributions, bug reports
  and feature requests.
- Others have also taken an interest in [maintaining their own
  forks](https://github.com/atc0005/go-teams-notify/network/members) of the
  original project. See those forks for other ideas/changes that you may find
  useful.

For more details, see the
[Releases](https://github.com/atc0005/go-teams-notify/releases) section or our
[Changelog](https://github.com/atc0005/go-teams-notify/blob/master/CHANGELOG.md).

## Supported Releases

| Series   | Example          | Status              |
| -------- | ---------------- | ------------------- |
| `v1.x.x` | `v1.3.1`         | Not Supported (EOL) |
| `v2.x.x` | `v2.6.0`         | Supported           |
| `v3.x.x` | `v3.0.0-alpha.1` | TBD                 |

The current plan is to continue extending the v2 branch with new functionality
while retaining backwards compatibility. Any breakage in compatibility for the
v2 series is considered a bug (please report it).

Long-term, the goal is to learn from missteps made in current releases and
correct as many as possible for a future v3 series.

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Usage

### Add this project as a dependency

See the [Examples](#examples) section for more details.

### Webhook URLs

#### Expected format

Valid webhook URLs for Microsoft Teams use one of several (confirmed) FQDNs
patterns:

- `outlook.office.com`
- `outlook.office365.com`
- `*.webhook.office.com`
  - e.g., `example.webhook.office.com`

Using a webhook URL with any of these FQDN patterns appears to give identical
results.

Here are complete, equivalent example webhook URLs from Microsoft's
documentation using the FQDNs above:

- <https://outlook.office.com/webhook/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c>
- <https://outlook.office365.com/webhook/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c>
- <https://example.webhook.office.com/webhookb2/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c>
  - note the `webhookb2` sub-URI specific to this FQDN pattern

All of these patterns when provided to this library should pass the default
validation applied. See the example further down for the option of disabling
webhook URL validation entirely.

#### How to create a webhook URL (Connector)

1. Open Microsoft Teams
1. Navigate to the channel where you wish to receive incoming messages from
   this application
1. Select `⋯` next to the channel name and then choose Connectors.
1. Scroll through the list of Connectors to Incoming Webhook, and choose Add.
1. Enter a name for the webhook, upload an image to associate with data from
   the webhook, and choose Create.
1. Copy the webhook URL to the clipboard and save it. You'll need the webhook
   URL for sending information to Microsoft Teams.
   - NOTE: While you can create another easily enough, you should treat this
     webhook URL as sensitive information as anyone with this unique URL is
     able to send messages (without authentication) into the associated
     channel.
1. Choose Done.

Credit:
[docs.microsoft.com](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using#setting-up-a-custom-incoming-webhook),
[gist comment from
shadabacc3934](https://gist.github.com/chusiang/895f6406fbf9285c58ad0a3ace13d025#gistcomment-3562501)

### Examples

#### Basic

This is an example of a simple client application which uses this library.

- `Adaptive Card`
  - File: [basic](./examples/adaptivecard/basic/main.go)
- `MessageCard`
  - File: [basic](./examples/messagecard/basic/main.go)

#### Specify proxy server

This is an example of a simple client application which uses this library to
route a generated message through a specified proxy server.

- `Adaptive Card`
  - File: [basic](./examples/adaptivecard/proxy/main.go)
- `MessageCard`
  - File: [basic](./examples/messagecard/proxy/main.go)

#### User Mention

These examples illustrates the use of one or more user mentions. This feature
is not available in the legacy `MessageCard` card format.

- File: [user-mention-single](./examples/adaptivecard/user-mention-single/main.go)
- File: [user-mention-multiple](./examples/adaptivecard/user-mention-multiple/main.go)
- File: [user-mention-verbose](./examples/adaptivecard/user-mention-verbose/main.go)
  - this example does not necessarily reflect an optimal implementation

#### Tables

These examples illustrates the use of a [`Table`][adaptivecard-table]. This
feature is not available in the legacy `MessageCard` card format.

- File: [table-manually-created](./examples/adaptivecard/table-manually-created/main.go)
- File: [table-unordered-grid](./examples/adaptivecard/table-unordered-grid/main.go)
- File: [table-with-headers](./examples/adaptivecard/table-with-headers/main.go)

#### Set custom user agent

This example illustrates setting a custom user agent.

- `Adaptive Card`
  - File: [custom-user-agent](./examples/adaptivecard/custom-user-agent/main.go)
- `MessageCard`
  - File: [custom-user-agent](./examples/messagecard/custom-user-agent/main.go)

#### Add an Action

This example illustrates adding an [`OpenUri`][msgcard-ref-actions]
(`MessageCard`) or [`OpenUrl`][adaptivecard-ref-actions] Action. When used,
this action triggers opening a URL in a separate browser or application.

- `Adaptive Card`
  - File: [actions](./examples/adaptivecard/actions/main.go)
- `MessageCard`
  - File: [actions](./examples/messagecard/actions/main.go)

#### Toggle visibility

These examples illustrates using
[`ToggleVisibility`][adaptivecard-ref-actions] Actions to control the
visibility of various Elements of an `Adaptive Card` message.

- File: [toggle-visibility-single-button](./examples/adaptivecard/toggle-visibility-single-button/main.go)
- File: [toggle-visibility-multiple-buttons](./examples/adaptivecard/toggle-visibility-multiple-buttons/main.go)
- File: [toggle-visibility-column-action](./examples/adaptivecard/toggle-visibility-column-action/main.go)
- File: [toggle-visibility-container-action](./examples/adaptivecard/toggle-visibility-container-action/main.go)

#### Disable webhook URL prefix validation

This example disables the validation webhook URLs, including the validation of
known prefixes so that custom/private webhook URL endpoints can be used (e.g.,
testing purposes).

- `Adaptive Card`
  - File: [disable-validation](./examples/adaptivecard/disable-validation/main.go)
- `MessageCard`
  - File: [disable-validation](./examples/messagecard/disable-validation/main.go)

#### Enable custom patterns' validation

This example demonstrates how to enable custom validation patterns for webhook
URLs.

- `Adaptive Card`
  - File: [custom-validation](./examples/adaptivecard/custom-validation/main.go)
- `MessageCard`
  - File: [custom-validation](./examples/messagecard/custom-validation/main.go)

## Used by

See the Known importers lists below for a dynamically updated list of projects
using either this library or the original project.

- [this fork](https://pkg.go.dev/github.com/atc0005/go-teams-notify/v2?tab=importedby)
- [original project](https://pkg.go.dev/github.com/dasrick/go-teams-notify/v2?tab=importedby)

## References

- [Original project](https://github.com/dasrick/go-teams-notify)
- [Forks of original project](https://github.com/atc0005/go-teams-notify/network/members)

- Microsoft Teams
  - MS Teams - adaptive cards
  ([de-de](https://docs.microsoft.com/de-de/outlook/actionable-messages/adaptive-card),
  [en-us](https://docs.microsoft.com/en-us/outlook/actionable-messages/adaptive-card))
  - MS Teams - send via connectors
  ([de-de](https://docs.microsoft.com/de-de/outlook/actionable-messages/send-via-connectors),
  [en-us](https://docs.microsoft.com/en-us/outlook/actionable-messages/send-via-connectors))
  - [adaptivecards.io](https://adaptivecards.io/designer)
  - [Legacy actionable message card reference][msgcard-ref]

[githubtag-image]: https://img.shields.io/github/release/atc0005/go-teams-notify.svg?style=flat
[githubtag-url]: https://github.com/atc0005/go-teams-notify

[goref-image]: https://pkg.go.dev/badge/github.com/atc0005/go-teams-notify/v2.svg
[goref-url]: https://pkg.go.dev/github.com/atc0005/go-teams-notify/v2

[license-image]: https://img.shields.io/github/license/atc0005/go-teams-notify.svg?style=flat
[license-url]: https://github.com/atc0005/go-teams-notify/blob/master/LICENSE

[msgcard-ref]: <https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference>
[msgcard-ref-actions]: <https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#actions>

[adaptivecard-ref]: <https://adaptivecards.io/explorer>
[adaptivecard-ref-actions]: <https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/getting-started>
[adaptivecard-user-mentions]: <https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#mention-support-within-adaptive-cards>
[adaptivecard-table]: <https://adaptivecards.io/explorer/Table.html>

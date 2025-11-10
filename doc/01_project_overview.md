
# Test-Flow - Project Overview

**Project Name:** Test-Flow  
**Type:** Practice Project / Portfolio Piece / SaaS Tool  

## Overview

Test-Flow is a platform designed to empower **non-technical users**—like Product Managers, Designers, and Support Analysts—to create automated end-to-end (E2E) tests, UI validations, and basic performance checks using **natural language inputs**. It removes the dependency on software developers for creating automated tests, making testing accessible, fast, and reliable.

The platform is inspired by real-world problems where small to medium companies struggle to maintain QA resources. For example, teams manually testing message flows in a chat application or validating UI against Figma designs.  

## Core MVP Features

- **Natural Language Test Creation**: Users describe tests in English.  
  - Example: `"Check if clicking the Send button in a chat app sends the message and displays it in the conversation."`
- **Automatic Test Generation**: AI generates executable **Playwright** scripts.
- **App-Agnostic Design**: Works with any web application using generic selectors.
- **Reporting**: Human-readable pass/fail results, execution times, and error details.
- **Configurable**: Optional `config.yaml` for custom selectors and overrides.
- **CLI Interface**: MVP is interactive CLI using **Go**.

## Long-Term Vision

- SaaS platform with web frontend, cloud execution, dashboards, RAG, and user memory.
- Open-source CLI for adoption and learning.
- Paid SaaS tiers with enterprise features, Figma integration, multi-user support, and optional dedicated AI instances.

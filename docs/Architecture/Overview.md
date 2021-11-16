## Overview

There are three primary components in the Switchboard architecture:
1. [[Client]] -- parses and renders [[Template Reference|templates]], which is the primary interface by which developers interact with the platform.
2. [[Worker]] -- deploys [[Resource Reference|resources]] using prebuilt [[Driver Interface|drivers]].
3. API -- the Porter API server.

When a developer interacts with the dashboard or CLI to perform an action, the client will construct a `ResourceGroup`, which subsequently gets applied by the worker. The dashboard/CLI also read data from the API server. 

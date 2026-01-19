# Update on AI Agent Forecasting
At the end of September 2025, I [wrote about my experiment](https://www.samuelsforecasts.com/blog/ai-agent-forecasting) with having AI agents forecast on this website in the same way I do. Now a couple of months have passed and the agents have been able to complete some forecasts, so I thought I would share some learnings and observations.

## Performance
The performance of all Agents was stellar as you can see in the table below. Gemini even beat me (which I guess says a lot about my forecasting ability hehe). Here's how the agents performed:

| forecaster| avg_brier  | forecasts |
| --------- | ---------- | --------- |
| Gemini    | 0.1112     | 30        |
| Samuel    | 0.1143     | 38        |
| Grok      | 0.1146     | 31        |
| Condorcet | 0.1252     | 35        |
| GPT-5     | 0.1283     | 24        |
| Opus      | 0.1677     | 29        |

The spread is not that big between the agents, apart from Opus which is an outlier due to being quite confidently wrong on a couple of forecasts (see e.g [this one](https://www.samuelsforecasts.com/forecast/237)). Unfortunately, the performance of each agent is correlated with how often I ran the agent, which in itself is strongly correlated with the cost per million tokens for the model. Gemini and Grok were the cheapest to run, hence they had more opportunities to update their forecasts and thus it is expected that they would perform slightly better. (This pattern holds for myself as well, this fall I averaged close to 1 forecast per question. Much lower than Gemini's 2.8 forecasts per question.)

The big counterexample of more forecast = better performance is the multi-agent system Condorcet, which has forecasted most of all agents. Yet it does not perform very well. Part of the reason for this was a bug in my tools and another is that the agent struggled to maintain coherence over longer time-periods and multiple sub-agent handoffs. I'll discuss this a bit more below.

## Common issues
Running AI agents autonomously for weeks revealed some more or less predictable failure modes. 

### Misunderstanding the resolution criteria
Even though the agents have access to a tool that gives them the resolution criteria, they often misunderstood the conditions for resolution. One example that came up multiple times were questions about interest rates and CPI, often these questions are similar to "will CPI get below X% before Y date?". Quite often models failed to take the date the forecast was created when forecasting these questions, so if CPI had been below X% before the question was even created, they would forecast close to 100%. 

This is partly a misunderstanding on the agent's part, but it is also a failure of writing resolution criteria that are unambiguous. This is something that I'm going to work on going forward.

### Coherence over longer periods of time
In some cases, the agents ran for extended periods of time. Some runs went on for more than 20 minutes. During this time, the agent executes lots of tools, such as finding new forecasts, fetching information, and creating forecasts. Over time, it becomes more and more difficult for the agent to keep track of its goal. This is especially true for Condorcet, the multi-agent system. 

Due to a combination of a poorly designed and implemented communication system and the inability for some models to keep on track and adhere to _all_ instructions, the multi-agent system quite often lost track of which tasks it had completed. 

This led to the agent creating multiple forecasts for one question within just a couple of minutes of each other, sometimes with quite different probabilities.

### Over-forecasting
Like many human forecasters, the agents tend to forecast too often and they tend to update too much on small information changes. The prompt I used had no guidelines on when the agents should update a forecast and not, which often led to agents updating forecasts by small or no amounts. 

Condorcet also often swings wildly in its forecasts, despite the underlying information not changing much. Here's a particularly striking [example](https://www.samuelsforecasts.com/forecast/257) where the agent forecasts a low probability at first, but then it changes it's mind to near certainty after a statistical analysis. Just a week later, it reduces the probability to less than 10% again despite nothing changing. 

This is partly due to the prompt not being very opinonated on how the agent should do forecasting. It is free to choose the workflow it believes is best and it can forecast whatever it wants. This freedom is interesting in an experiment, but in a setting where performance was required, it would likely require stricter guidelines.

## Going Forward
As I've already hinted above, there are lots of things that could be improved with these agents and the architectures. I've also realized that I'm a little bit hindered by how much I'm willing to spend on tokens. Unless someone wants to sponsor my API bills (my DMs are open), I'm going to focus on only two agents going forward: one single-model agent that only uses the current SOTA model and Condorcet, the multi-agent system I've designed. This will allow me to focus on the architectures to improve performance.

More specifically, I want to explore ways to improve coherence over longer periods of time and longer workflows. My current ideas to do this are:
- give agents access to a persistent filesystem which can make context close to infinite with clever context management and exploration tools.
- improve management and "oversight" of subagents for the Orchestrator agent
- give agents tools to crate their own data sources. Forecasting is heavily influenced by the available data and just having Perplexity as a data source is not good enough, especially if I want agents to rely on statistical methods as well as qualitative forecasting.
- work on methods to rapidly test and improve prompts. It is clear that prompting is very important for good and consistent performance and I haven't really experimented much with the prompts so far.

I'll continue running these experiments and will share another update in a few months. In the meantime, you can see the agents' live forecast performance on the [website](https://www.samuelsforecasts.com).

## AI Agent Forecasting
This fall I've been working on an experiment to learn more about AI agents and to benchmark AI models in my own way. I've built two different AI agent architectures for forecasting, one single-agent architecture and one multi-agent system. In this post I'll introduce the AI agents and their architecture.

### Forecasting as a benchmark for AI models
Before going into the specifics of the implementation, let's think a little bit about why forecasting is a good domain for AI agents (setting aside the fact that it is a domain I'm interested in) and why it is suitable for evaluating AI model capabilities.

Forecasting requires gathering evidence, figuring out how to evaluate evidence, and building a model of the different outcomes. To perform well, an AI model needs to not only be intelligent, it also needs to have access to good tools that allow it to gather information, which it can later use. It also requires decision-making capabilities, as the model/agent needs to be able to decide when it needs more information and when it is ready to commit to a final forecast. The AI model also needs to be creative in how it uses the tools and information it gathers to create a calibrated and well-reasoned forecast.

All of these factors, along with me already having this website and some pre-made tools for forecasting, make forecasting a good testbed for how good different LLMs are at reasoning and decision-making.

### Single-agent forecasting
I started out by building a single-agent architecture from Anthropic's [Claude Quickstarts](https://github.com/anthropics/claude-quickstarts). To be able to use any model, I refactored that setup into using the OpenAI sdk. This makes it possible to use OpenRouter to select any model through their API.

To make forecasting possible for the agents, I created the following tools to interact with my website:
- `get_forecasts()`: fetches open forecast questions from this website that the agent can forecast on.
- `get_forecast_data(forecast_id)`: fetches detailed information, like resolution criteria and date the forecast was created for a specific forecast_id.
- ` get_forecast_points(forecast_id)`: fetches all previous forecasts made by the agent on a specific question.
- `update_forecast(forecast_id, probability, reason)`: creates a new forecast for a question with reasoning.

To allow the agents to fetch information from the internet, I created a `query_perplexity(query)` tool which sends the agents query to Perplexity and returns a response for the agent to ingest. I've also created a `request_feedback()` tool which allows the AI agent to request feedback from an expert forecaster (GPT-5 with a specific critic prompt). 

These tools all combine into a linear flow in which the agents call tools in a fairly linear order. From fetching open forecast questions, to selecting one or a couple questions to focus on, to getting the details of the questions and the previous forecasts, and then searching for information using Perplexity to synthesize and update the question with a new forecast. 

ADD A PICTURE OF THE FLOW HERE. 

Here's some logs for a Gemini 2.5 Pro forecast:
```
============================================================
NEW FORECASTING SESSION: single_agent_20250928_195427
Started: 2025-09-28 19:54:27
============================================================
19:54:27 | [INFO]  [ORCH] [System] Single Agent Test Starting
19:54:27 | [INFO]  [ORCH] [ForecastingAgent] Received task | Go ahead and forecast!
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:30 | [INFO]  [LLM] google/gemini-2.5-pro | 1,975 tokens
19:54:30 | [INFO]   ├─ [REASONING]
         |     │ **Listing Available Forecasts**
         |     │ 
         |     │ I'm currently working on identifying the available forecasts. I've created a plan to get the list of forecasts, and I am now gathering data. It is a simple process and the list of forecast options should be ready shortly.
         |     │ 
         |     │ 
         |     │ 
19:54:30 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecasts
HTTP Request: GET http://forecasting-389105.ey.r.appspot.com/forecasts/stale-and-new/21 "HTTP/1.1 200 OK"
19:54:30 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecasts ✓ [{'id': 253, 'question': 'Will Arsenal win the 2025-2026 Premier League season?', 'category': 'Sports', 'created': '2025-08-26T14:31:18.845727Z', 'use...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:37 | [INFO]  [LLM] google/gemini-2.5-pro | 6,591 tokens
19:54:37 | [INFO]   ├─ [REASONING]
         |     │ **Prioritizing Inquiry**
         |     │ 
         |     │ I'm currently focusing on identifying a manageable set of questions. I'm leaning towards those with clear answers, tight deadlines, and accessible, dependable data. Right now, I'm considering the likelihood of a new ceasefire between Israel and Hamas before the end of the year.
         |     │ 
         |     │ 
         |     │ **Detailing My Questions**
         |     │ 
         |     │ I've just refined the question parameters; resolution criteria, timeframe, and data sources are now well-defined. Specifically, I am now digging deep into the nuances of question 238 - focusing on the Israel-Hamas ceasefire.
         |     │ 
         |     │ 
         |     │ 
19:54:37 | [INFO]   ├─ [RESPONSE]
         |     │ I will begin by selecting a few promising questions to investigate, prioritizing those with clear resolution criteria, a relatively near-term resolution date, and a high likelihood of available public data.
         |     │ 
         |     │ My initial shortlist is:
         |     │ 1. **238: Will there be a new ceasefire between Israel and Hamas before 2026?** (High-impact, time-sensitive, good information flow)
         |     │ 2. **228: Will the US enter a recession in 2025?** (High-impact, medium-term, data-rich)
         |     │ 3. **220: Will there be a case of human-to-human transmission of H5N1 in 2025?** (High-impact, reliant on official health reports)
         |     │ 4. **242: Will the US attack Iran a second time in 2025?** (High-impact, time-sensitive, news-driven)
         |     │ 5. **202: Will there be a ceasefire in the Russia-Ukraine conflict before 2026?** (High-impact, ongoing, extensive news coverage)
         |     │ 
         |     │ I will start with the first question: **238: Will there be a new ceasefire between Israel and Hamas before 2026?**
         |     │ 
         |     │ First, I will retrieve the full details for this question.
19:54:37 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_data(forecast_id: 238)
HTTP Request: GET http://forecasting-389105.ey.r.appspot.com/forecasts/238 "HTTP/1.1 200 OK"
19:54:37 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_data ✓ {'success': True, 'data': {'id': 238, 'question': 'Will there be a new ceasefire between Israel and Hamas before 2026?', 'category': 'Conflict, World'...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:41 | [INFO]  [LLM] google/gemini-2.5-pro | 6,582 tokens
19:54:41 | [INFO]   ├─ [REASONING]
         |     │ **Assessing the Forecast Question**
         |     │ 
         |     │ I'm now diving into the forecasting task. My primary focus is identifying and understanding the specific criteria for resolution: a new, mutually agreed-upon ceasefire lasting over a month. I need to thoroughly analyze the nuances of this definition.
         |     │ 
         |     │ 
         |     │ **Framing the Core Question**
         |     │ 
         |     │ I'm currently mapping out the core parameters of this forecast. First, confirming the question's constraints, which is the timeline and conditions of the ceasefire agreement. Then, analyzing the specific definition for a successful forecast: a new, public ceasefire of at least 30 days agreed upon before 2026. This is a critical first step.
         |     │ 
         |     │ 
         |     │ 
19:54:41 | [INFO]   ├─ [RESPONSE]
         |     │ I will now check to see if I have previously forecasted on this question.
19:54:41 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_points(forecast_id: 238)
HTTP Request: POST http://forecasting-389105.ey.r.appspot.com/forecast-points/user "HTTP/1.1 200 OK"
19:54:41 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_points ✓ {'success': True, 'data': [{'id': 894, 'forecast_id': 238, 'point_forecast': 0.99, 'reason': "BLUF: I am at 0.99 that this has already resolved as 'Ye...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:46 | [INFO]  [LLM] google/gemini-2.5-pro | 7,095 tokens
19:54:46 | [INFO]   ├─ [REASONING]
         |     │ **Assessing the Initial Forecast**
         |     │ 
         |     │ I've reviewed the initial forecast on the ceasefire question and the relevant background. Now, I'm integrating more recent developments and refining my probability estimates. It's a dynamic situation, but my initial assessment seems to hold, with some potential minor adjustments needed.
         |     │ 
         |     │ 
         |     │ **Evaluating the January Ceasefire**
         |     │ 
         |     │ I'm gathering information about the January 2025 ceasefire. I'm using Perplexity to search for reports on its duration and any breaches. The key question is: did it hold for the necessary 30+ days to validate my initial 0.99 probability? My analysis will hinge on these findings.
         |     │ 
         |     │ 
         |     │ 
19:54:46 | [INFO]   ├─ [RESPONSE]
         |     │ I have already made a forecast on this question on September 7, 2025, with a probability of 0.99. My reasoning was that a ceasefire had been agreed upon in January 2025 for 42 days, which would satisfy the resolution criteria. I will now check if this ceasefire held and if there have been any new developments.
19:54:46 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity(query_text: Israel-Hamas ceasefire January 2025 success)
19:54:51 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity ✓ The Israel-Hamas ceasefire that began on January 19, 2025, achieved initial success by halting active hostilities, facilitating the return of Israeli ...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:54 | [INFO]  [LLM] google/gemini-2.5-pro | 7,432 tokens
19:54:54 | [INFO]   ├─ [RESPONSE]
         |     │ The Perplexity query confirms that a 42-day ceasefire began on January 19, 2025, and that it held, at least for the first phase. The resolution criteria for the forecast is "a new ceasefire is agreed upon by both Israel and Hamas AND the ceasefire holds for more than 30 days... before Jan 1st 2026". Since a 42-day ceasefire was agreed upon and held, the question should resolve as "Yes." My previous forecast of 0.99 was accurate. There is no need to update this forecast.
         |     │ 
         |     │ I will now move to the next question on my shortlist: **228: Will the US enter a recession in 2025?**
         |     │ 
         |     │ I will start by getting the full details for this question.
19:54:54 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_data(forecast_id: 228)
HTTP Request: GET http://forecasting-389105.ey.r.appspot.com/forecasts/228 "HTTP/1.1 200 OK"
19:54:54 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_data ✓ {'success': True, 'data': {'id': 228, 'question': 'Will the US enter a recession in 2025?', 'category': 'Finance, Economy', 'created': '2025-01-11T18:...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:56 | [INFO]  [LLM] google/gemini-2.5-pro | 7,564 tokens
19:54:56 | [INFO]   ├─ [RESPONSE]
         |     │ Now, I'll check my previous forecasts for this question.
19:54:56 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_points(forecast_id: 228)
HTTP Request: POST http://forecasting-389105.ey.r.appspot.com/forecast-points/user "HTTP/1.1 200 OK"
19:54:56 | [INFO]   ├─ [TOOL] [ForecastingAgent] → get_forecast_points ✓ {'success': True, 'data': [{'id': 891, 'forecast_id': 228, 'point_forecast': 0.25, 'reason': "BLUF: I am at 0.25 that the NBER will declare a US reces...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:54:58 | [INFO]  [LLM] google/gemini-2.5-pro | 8,178 tokens
19:54:58 | [INFO]   ├─ [RESPONSE]
         |     │ I previously forecasted on this question on September 7, 2025, with a probability of 0.25. My reasoning was based on the NBER's historical lag in declaring recessions, and conflicting economic forecasts. I will now gather fresh information to see if the outlook has changed. I'm interested in recent GDP reports, unemployment data, and statements from the NBER or prominent economists.
19:54:58 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity(query_text: US economic outlook September 2025)
19:55:07 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity ✓ The **US economic outlook for September 2025** indicates **below-trend but positive growth**, with slowing momentum influenced by tariff effects, pers...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:55:11 | [INFO]  [LLM] google/gemini-2.5-pro | 9,065 tokens
19:55:11 | [INFO]   ├─ [RESPONSE]
         |     │ The latest information from September 2025 points to a slowdown, but not a recession. GDP growth is expected to be positive, albeit at a lower rate. The labor market is weakening, but not collapsing. Inflation remains a concern. The consensus forecast seems to be for continued, slow growth.
         |     │ 
         |     │ The question resolves on an NBER declaration of a recession *before January 1st, 2026*. The NBER's "three Ds" for a recession are that it must be deep, diffuse, and durable. Given that we are already in late September 2025, a recession would have to begin very soon and be very severe for the NBER committee to declare it so quickly. The NBER's own FAQ states that they often wait several months to a year to declare a recession, to be sure of the data.
         |     │ 
         |     │ Considering the current data (slow growth, not contraction), the high bar for an NBER declaration, and the short remaining time, the probability of a "Yes" resolution seems even lower than my previous 0.25 estimate. I will lower my forecast. I want to search for the historical lag time of NBER recession announcements.
19:55:11 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity(query_text: NBER recession announcement lag time)
19:55:17 | [INFO]   ├─ [TOOL] [ForecastingAgent] → query_perplexity ✓ The National Bureau of Economic Research (NBER) typically announces the start and end dates of recessions with a significant lag, often many months or...
HTTP Request: POST https://openrouter.ai/api/v1/chat/completions "HTTP/1.1 200 OK"
19:55:26 | [INFO]  [LLM] google/gemini-2.5-pro | 10,033 tokens
19:55:26 | [INFO]   ├─ [RESPONSE]
         |     │ The information about the NBER's lag time is very helpful. While the 2020 recession announcement was unusually fast (4 months), other recessions have had much longer announcement lags (e.g., the 2001 recession announcement took 20 months). The NBER's process emphasizes accuracy, so they are generally slow to make a call.
         |     │ 
         |     │ Given that it's the end of September 2025, and there's no clear evidence of a sharp, deep, and broad-based downturn, the likelihood of an NBER declaration within the next three months is very low. Even if a recession were to start *today*, the NBER would be unlikely to declare it before the end of the year. The economic data, while showing a slowdown, doesn't point to the kind of severe contraction that would prompt a rapid NBER announcement.
         |     │ 
         |     │ Therefore, I'm confident in significantly lowering my forecast. A probability of 0.05 seems more appropriate, reflecting the very low chance of an NBER declaration in the required timeframe, but acknowledging the non-zero possibility of a sudden and severe economic shock.
         |     │ 
         |     │ I will now update my forecast.
19:55:26 | [INFO]   ├─ [TOOL] [ForecastingAgent] → update_forecast(point_forecast: 0.05, reason: 'BLUF: I am lowering my forecast to 0.05. The NB...', forecast_id: 228)
HTTP Request: POST http://forecasting-389105.ey.r.appspot.com/users/login "HTTP/1.1 200 OK"
HTTP Request: POST http://forecasting-389105.ey.r.appspot.com/api/forecast-points "HTTP/1.1 201 Created"
19:55:26 | [INFO]   ├─ [TOOL] [ForecastingAgent] → update_forecast ✓ {'success': True, 'data': 'Forecast point created'}
```

The linear and not-so-creative flow is mainly due to the (prompt)[https://github.com/bayesmaxxing/agent-forecaster/blob/main/prompt.md], which is quite instructive and it does not open up for creative interpretation of the forecasting process. Overall, the limited list of tools and the rigid prompt creates a limiting factor for the creativeness of the agents, there's just not that many different ways they can forecast in this architecture. To deal with this, I built the multi-agent forecasting system.

### Multi-agent forecasting system
To allow for more creativeness and better forecasts, I wanted to build a system that could make itself as complex as it wished. Hence I built two new important tools that require some explanation.

#### Subagent Manager Tool






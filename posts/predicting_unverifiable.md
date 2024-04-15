Over the past years, I have on a few occasions been writing about prediction markets. I’ve shown how they are valuable in both business and politics and I’ve discussed Robin Hanson’s suggestion of incorporating prediction markets into democracy—what he calls Futarchy. But when discussing prediction markets with people, one of the most common questions is how one would create prediction markets on questions that have outcomes that are either unverifiable or, at least, very hard to verify. This is a relevant question, especially in politics where a lot of the questions we care about are hard to verify.

In many questions, it is possible to identify some identifiers of success from some policy. If you suggest some school reform, which aims to improve performance, then one of the identifiers of success could be performance on international benchmarks such as PISA—but how would you unambiguously decide whether that was casually due to the school reform rather than due to some other policy choice, or pure randomness? 

In these types of situations, we would still want to gain trustworthy information about how likely it is that our policies are. (Though the fact that politicians are unwilling to try prediction markets suggests that they may not be very interested in knowing how likely their favorite policies are to be successful.) Regular prediction markets are not very well equipped to handle ambiguous termination, they require well-defined criteria to be able to elicit trustworthy and incentive-compatible predictions. To deal with this, we need some other sort of mechanism and one such mechanism has been proposed by a recent paper. (https://arxiv.org/pdf/2306.04305.pdf) 

The paper suggests a sequential peer prediction mechanism, that they call a prediction market. But unlike prediction markets, this mechanism does not really incorporate tradable securities in the way that prediction markets do. Instead, this mechanism is based on sequential elicitation of predictions, where each agent participating will compare his belief to a reference agent. There are multiple ways of setting up such a market, and the one that the authors prefer, due to the comparative simplicity, is to set up a market that has a small probability, *alpha*, of terminating after each prediction. The agent that happens to set the final prediction, we will call this agent *T*, will be the reference agent for one subset of all the predictors. The reference agent *T* is used to determine the payoff for the initial *T-k* predictors using negative cross-entropy. The remaining *k* agents will be paid a flat rate *R*. 

This setup, with *k* agents that the reference agent can observe, and *T-k* agents to get the market going, means that it is possible, with *k* large enough and *alpha* small enough, to make a prediction mechanism like this incentive-compatible within a small enough bound. One cannot completely eliminate a very small incentive to deviate from truthful reporting of one’s beliefs, but a larger *k* makes that incentive small enough to not matter. 

After each prediction, the market continues with probability 1-alpha. Setting *alpha = 1/(T+k)* would let you on average elicit *T+k* predictions. This setup leads to a self-resolving prediction market, meaning that you don’t need to wait for measurements or experiments to know how to resolve the market. The market resolves without any real answer to the question, yet each agent has a strong incentive to report his true beliefs. 

Compared to the sort of conditional prediction markets that I, based on Robin Hanson’s ideas, have argued for, these markets are potentially cheaper in that they do not need any verification to be carried out. A conditional market asks questions of the kind “conditional on policy X being enacted, will outcome Y improve?” and thus to get a resolution, you must be able to establish that the outcome Y did improve as a result of policy X. To do this can be expensive, either you need some sort of experiment or you need some other way of identifying causality. 

With this self-resolving market, you do not need that for payouts. Of course, there are other problems. Most notably, this mechanism does not really take into account multiple predictions in the way a market mechanism does. Further, the key result of the paper is that the truthfulness of the reported beliefs relies on the reference agent having access to the k agents before him that have reported their conditionally independent private signals. This assumption about informational substitutes is crucial, if this does not hold, which it doesn’t if the agents have a cost of exerting effort or in gathering information. So this peer mechanism is not really fit to deal with new information becoming available to agents or when the private signals to the participating agents are not informational substitutes. 

Yet, it seems to me like one could use this mechanism as another tool in the toolbox for policy makers to improve their knowledge of the effects of their policies. This mechanism seems useful for unverifiable questions that need to be resolved comparatively quickly, where you don’t want to incentivize agents to gather more information between their predictions and where it seems likely that the assumption of informational substitutes holds. (I will have to think more about questions where I think this type of assumption is reasonable.) On the other hand, conditional prediction markets are more useful in cases where the outcomes are more easily specified and measured. 

Having such a toolbox, combined with policy experiments, would be useful for politicians and policy makers that are concerned with extracting the most out of their policies, or at least reducing the likelihood of enacting policies that are net negative on society.  
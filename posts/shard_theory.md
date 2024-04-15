Human values are elusive, they are hard to completely specify, they seem to be very contextually dependent, and they seem to arise from a mix of genetics and environmental exposure. Thus it is not very surprising that few gears-level accounts of how humans form their values exist. Most of what can be found online are handwavy and vague accounts of how values form from culture and socialization—they rarely mention how this process happens in the brain. In this post, I want to introduce shard theory, the theory of value formation proposed by Alex Turner and Quintin Pope. (https://www.alignmentforum.org/posts/iCfdcxiyr2Kj8m8mT/the-shard-theory-of-human-values)

First, we need to specify what we mean by values. In this post, we will take it to mean “a contextual influence on decision-making”, a fairly broad definition. In this sense, a truthfulness-value will influence me to make decisions that are more truthful in most situations. (Unless perhaps some other value overrules it.) We also need a few base assumptions for the theory to make sense. *(1) The Cortex is basically randomly initialized, (2) the brain does self-supervised learning and (3) the brain does reinforcement learning.* 

**Shard Theory**

Now we have the assumptions as well as the definition of a value. The simplified version of the theory goes:

DNA → 
Hard-coded reward circuitry in the brain → 
Obtain reward through interaction with the environment → 
Reinforcement learning/credit assignment on the computations that led to reward → 
Reinforced computations become heuristics, or shards → 
The shards and the learned world model are intertwined → 
The agent learns a crude planning algorithm → 
the agent becomes a planner → 
Shards are used in decision-making as heuristics that lead to reward → 
Shards “bid” for their preferred plan, leading to contextual influence on decision-making. 

Let’s go a bit more into detail. DNA is widely believed to specify some reward function, that helps with survival, and this is initialized in the brain sometime in-utero. At the beginning of life, the agent only has the crude reward circuitry to guide it to better and better actions. Turner and Pope explain it like this:

‘... suppose that the genome specifies a reward circuit which takes as input the state of the taste buds and the person’s metabolic needs, and produces a reward if the taste buds indicate the presence of sugar while the person is hungry. By assumption 3 in section I, the brain does reinforcement learning and credit assignment to reinforce circuits and computations which led to reward. For example, if a baby picks up a pouch of apple juice and sips some, that leads to sugar-reward. The reward makes the baby more likely to pick up apple juice in similar situations in the future. 

Therefore, a baby may learn to sip apple juice which is already within easy reach. However, without a world model (much less a planning process), the baby cannot learn multi-step plans to grab and sip juice. If the baby doesn’t have a world model, then she won’t be able to act differently in situations where there is or is not juice behind her. Therefore, the baby develops a set of shallow situational heuristics which involve sensory preconditions like “IF juice pouch detected in center of visual field, THEN move arm towards pouch.” The baby is basically a trained reflex agent.’ 

Of course, this type of proto-planning heuristic is unlikely to lead to more complex behaviors and values unless there’s a model of the world. Executing the heuristics leads to learning to associate proto-planning heuristics with concepts in the world model, in the juice-shard example, the intertwining with the world model happens when the baby sees that the juice often stands in the kitchen when it sees the pouch. Later on, the child sees its parent get the pouch of juice from the refrigerator, and thus the brain is reinforcing the fridge to be associated with the juice-shard. 

From the continued activations of these proto-planning shards, the brain now learns a crude planning algorithm using the world model and its shards. It is now that the baby becomes a planner that can form coherent plans to get things it values. For plans to be coherent, the baby has to weigh different shards against each other in a process that can be seen as a kind of “bidding” process. Where different shards are bidding for plans that lead to their outcome, i.e the juice shard bids for plans that lead to juice consumption. 

In this way, we have a process that leads to values according to the above definition. And humans, according to the theory, are not reward-maximizers but rather shard-executors. The important feature of this theory is that the values are shaped by the crude reward system without actually binding to the activations of that system. This reduces the incentives for wireheading and can explain why humans generally don’t want value drift. 

With the theory, it is possible to retrodict many human behaviors such as different biases. (Read the original post for an in-depth treatment of human behavior in light of shard theory.)

**Alignment of Artificial Intelligences**

Having a gears-level theory of human value formation seems like an important part of understanding why humans are surprisingly aligned. The only data points we have in the alignment of intelligent systems are humans, and we are very aligned on the whole. (Of course, there are some salient examples of non-aligned humans. But they are not the rule.) 

Based on the theory, the authors conjecture that humans are aligned, or that they convergently care about certain objects in the real world, because of the order in which the brain learns abstractions. 

To me, it seems like there are two main, though not mutually exclusive, possibilities from this point. (1) The learning algorithm, the combination of self-supervised learning and reinforcement learning, that drives the initial learning is responsible for the order in which the brain learns abstractions or (2) another biological process, that also is hard-coded in the brain, is responsible for the order in which the brain learns abstractions.

In case (1), it seems like implementing the learning algorithm and crude reward system that humans are hard-coded within an AI could lead to the AI also developing shards of value and coming to value certain objects in the real world. Then we are back to one of the original problems of AI alignment, specifying rewards that lead to values that are aligned with human values. Only this time, we need to worry less about the consequences of optimizing directly on reward and more about what kinds of shards we can get out of the reward function. (The intuition here is that it seems unlikely that a reward function that rewards calorie intake can lead to valuing future generations of humans.) In the other case (2), it seems like it would be harder to implement shard theory to solve alignment. Though it may be that we can model the biological processes well enough in the future to implement it.  

One thing that both (1) and (2) seem to imply is that training will be important. If shards are heavily influenced by the order and number of times computations happen, then that high path dependency will mean that designing training processes will be important. Overall, evidence of humans shows that there’s a fair bit of path dependency in human value formation. In general, I believe that a high path dependency world means that we will have to be extra careful about how we design the reward system, learning algorithm, and training process. In a low path dependency world, small perturbations in the initial conditions and training process are less likely to lead to wildly different generalizations. Thus, showing that shard theory would be implementable in artificial systems would be a step forward. (Current ML systems seem less path dependent than shard theory would predict, which means that current algorithms and architectures do not create shards in the way the brain does.)

This ties into a longer discussion about low and high path dependence in the current ML regime and the limit of training processes, which I hope to write more about in the future. I think that relying on creating a good training process is too dangerous to rely on for AI alignment, so in that case, I hope that I’m wrong about shard theory needing close to perfect training to create aligned AI. With that said, there’s still much to learn about shard theory and its potential implications for alignment and I’m looking forward to reading more from Alex, Quintin, and other researchers working on it.

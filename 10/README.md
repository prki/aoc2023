AOC 2023 10
===

Due to lack of time and other priorities, part 2 was not implemented.

Part 1 was implemented as a simple BFS over parsing into a proper graph structure.

Part 2's problem is to discover points enclosed by the pipes - with the pipes considered with a real "width". What that means is that the animal can squeeze e.g. between pipes ||.

In order to solve this, a proposed solution could be to create a polygon shape from the discovered graph. All possible points which may and may not be enclosed would be points not in the graph.
Then, each point would be tested to check whether the point is in polygon. Various algorithms exist, such as drawing a line from a point in any direction and checking count of intersections
with the polygon. In case it's even (line enters, line leaves), the point is not in the polygon, if it's odd, it is in the polygon.

The least painful way would probably be to find some math/game library which already contains some capabilities and just be careful with sizing the polygon and the points, but other than
that, there shouldn't be much more to the problem.

It is unclear if I will get back to this or not, since there are personal projects I need to complete by DEC 2023. If so, I should update this, hopefully.

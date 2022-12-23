package blueprint

type Blueprint struct {
	id          int
	robot_costs [][]int
}

type BlueprintStackItem struct {
	blueprint      *Blueprint
	curr_robots    []int
	curr_resources []int
	curr_time      int
}

func NewBlueprint(id int, robot_costs [][]int) *Blueprint {
	var B Blueprint
	B.id = id
	for i := range robot_costs {
		B.robot_costs = append(B.robot_costs, robot_costs[i])
	}

	return &B
}

func (B *Blueprint) GetID() int {
	return B.id
}

func (Bsi *BlueprintStackItem) copyBlueprintStackItem() *BlueprintStackItem {
	new_bsi := BlueprintStackItem{
		blueprint:      Bsi.blueprint,
		curr_robots:    make([]int, len(Bsi.curr_robots)),
		curr_resources: make([]int, len(Bsi.curr_resources)),
		curr_time:      Bsi.curr_time,
	}
	copy(new_bsi.curr_robots, Bsi.curr_robots)
	copy(new_bsi.curr_resources, Bsi.curr_resources)

	return &new_bsi
}

func (Bsi *BlueprintStackItem) getMostPossibleGeodes(max_time int) int {
	time_remaining := max_time - Bsi.curr_time
	curr_robot_production := time_remaining * Bsi.curr_robots[len(Bsi.curr_robots)-1]
	possible_robot_production := time_remaining * time_remaining / 2
	return curr_robot_production + possible_robot_production + Bsi.curr_resources[len(Bsi.curr_resources)-1]
}

func (Bsi *BlueprintStackItem) getTimeUntilRobot(robot_index int) (int, bool) {
	var robot_time int
	for i := range Bsi.blueprint.robot_costs[robot_index] {
		cost_diff := Bsi.blueprint.robot_costs[robot_index][i] - Bsi.curr_resources[i]
		if Bsi.curr_robots[i] == 0 && cost_diff <= 0 {
			continue
		} else if Bsi.curr_robots[i] == 0 {
			return 0, false
		}
		curr_robot_time := (cost_diff + Bsi.curr_robots[i] - 1) / Bsi.curr_robots[i]
		if curr_robot_time > robot_time {
			robot_time = curr_robot_time
		}
	}

	return robot_time + 1, true
}

func (Bsi *BlueprintStackItem) advanceTime(time int) {
	Bsi.curr_time += time
	for i := range Bsi.curr_robots {
		Bsi.curr_resources[i] += Bsi.curr_robots[i] * time
	}
}

func (Bsi *BlueprintStackItem) purchaseRobot(robot_index int) {
	Bsi.curr_robots[robot_index]++
	for i := range Bsi.blueprint.robot_costs[robot_index] {
		Bsi.curr_resources[i] -= Bsi.blueprint.robot_costs[robot_index][i]
	}
}

func (Bsi *BlueprintStackItem) getNextBlueprintStacks(max_time int) []*BlueprintStackItem {
	var next_blueprint_stack []*BlueprintStackItem

	// check which robot to build next
	best_time := max_time
	for i := len(Bsi.curr_robots) - 1; i >= 0; i-- {
		robot_time, ok := Bsi.getTimeUntilRobot(i)
		if !ok {
			continue
		}
		if robot_time <= 0 {
			robot_time = 1
		}
		if Bsi.curr_time+robot_time > max_time {
			continue
		}
		if robot_time > best_time { // don't build if takes longer than a geode robot
			continue
		} else if i == len(Bsi.curr_robots)-1 {
			best_time = robot_time
		}

		if i < len(Bsi.curr_robots)-1 { // don't build if robots produce more than we can consume
			if Bsi.curr_robots[i] >= Bsi.blueprint.getMaxResourceCost(i) {
				continue
			}
		}

		next_blueprint_stack_item := Bsi.copyBlueprintStackItem()
		next_blueprint_stack_item.advanceTime(robot_time)
		next_blueprint_stack_item.purchaseRobot(i)

		next_blueprint_stack = append([]*BlueprintStackItem{next_blueprint_stack_item}, next_blueprint_stack...)
	}

	// option for no decision, just passing time
	if Bsi.curr_time < max_time && len(next_blueprint_stack) == 0 {
		next_blueprint_stack_item := Bsi.copyBlueprintStackItem()
		next_blueprint_stack_item.advanceTime(max_time - Bsi.curr_time)
		next_blueprint_stack = append(next_blueprint_stack, next_blueprint_stack_item)
	}

	return next_blueprint_stack
}

func (Bsi *BlueprintStackItem) isWorseThan(other_bsi *BlueprintStackItem) bool {
	if Bsi.curr_time < other_bsi.curr_time {
		return false
	}
	for i := range Bsi.curr_robots {
		if Bsi.curr_robots[i] > other_bsi.curr_robots[i] {
			return false
		}
	}
	for i := range Bsi.curr_resources {
		if Bsi.curr_resources[i] > other_bsi.curr_resources[i] {
			return false
		}
	}
	return true
}

func (B *Blueprint) getMaxResourceCost(resource_index int) int {
	var max_cost int
	for i := range B.robot_costs {
		if B.robot_costs[i][resource_index] > max_cost {
			max_cost = B.robot_costs[i][resource_index]
		}
	}
	return max_cost
}

func (B *Blueprint) GetGeodeProduction(time int) int {
	var max_geodes int
	first_stack_item := BlueprintStackItem{
		blueprint:      B,
		curr_robots:    make([]int, len(B.robot_costs)),
		curr_resources: make([]int, len(B.robot_costs)),
	}
	first_stack_item.curr_robots[0] = 1
	blueprint_stack := []*BlueprintStackItem{&first_stack_item}

	best_stack_item_map := make(map[int]*BlueprintStackItem)

	for len(blueprint_stack) != 0 {
		curr_stack_item := blueprint_stack[len(blueprint_stack)-1]
		blueprint_stack = blueprint_stack[:len(blueprint_stack)-1]

		if curr_stack_item.curr_time > time {
			continue
		}
		if curr_stack_item.curr_resources[len(curr_stack_item.curr_resources)-1] > max_geodes {
			max_geodes = curr_stack_item.curr_resources[len(curr_stack_item.curr_resources)-1]
		}
		if curr_stack_item.curr_time == time {
			continue
		}
		if curr_stack_item.getMostPossibleGeodes(time) < max_geodes {
			continue
		}

		curr_best, ok := best_stack_item_map[curr_stack_item.curr_time]
		if !ok || curr_best.isWorseThan(curr_stack_item) {
			for i := curr_stack_item.curr_time; i <= time; i++ {
				curr_best, ok := best_stack_item_map[curr_stack_item.curr_time]
				if !ok || curr_best.isWorseThan(curr_stack_item) {
					best_stack_item_map[i] = curr_stack_item
				} else {
					break
				}
			}
		} else if curr_stack_item.isWorseThan(curr_best) {
			continue
		}

		blueprint_stack = append(blueprint_stack, curr_stack_item.getNextBlueprintStacks(time)...)
	}

	return max_geodes
}

import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

const animalFormat = (level, family = {}, id) => {
    const animal = family[id]
    if (!animal) { return { name: "Unknown", value: 1024/(2**level), gender: -1 } }
    if (level >= 4) { return { ...animal, value: 1024/(2**level) } }
    return { ...animal, children: Object.values(animal.parents || {}).map(id => animalFormat(level+1, family, id))  }
}

const AnimalRadialTree = ({ animal = null }) => {
    const d3Container = useRef(null);
    useEffect(
        () => {
            if (animal && d3Container.current) {
                const width = d3Container.current.getBoundingClientRect().width;
                const height = d3Container.current.getBoundingClientRect().height*1.2;
                const svg = d3.select(d3Container.current)
                    .attr("width", width*2)
                    .attr("height", height*2)
                    .attr("viewBox", [-width/2, -height/2, width, height])
                    .attr("style", "width: 100%; height: 100%;");

                const data = { ...animal, gender: 2, children: Object.values(animal.parents || {}).map(id => animalFormat(1, animal.family, id)) }

                const color = (v) => {
                    if (v?.gender == undefined) { return "gray" }
                    return v?.gender == 0 ? "blue" : "red"
                }
                const radius = 928 / 2;

                // Prepare the layout.
                const partition = data => d3.partition()
                    .size([2 * Math.PI, radius])
                    (d3.hierarchy(data)
                        .sum(d => d.value)
                        /*.sort((a, b) => b.value - a.value)*/);

                const arc = d3.arc()
                    .startAngle(d => d.x0)
                    .endAngle(d => d.x1 - 0.005)
                    // .padAngle(d => Math.min((d.x1 - d.x0) / 2, 0.005))
                    // .padAngle(d => Math.min((d.x1 - d.x0) / 2, 0.005))
                    // .padRadius(radius / 2)
                    .innerRadius(d => d.y0)
                    .outerRadius(d => d.y1 - 1);

                const root = partition(data);

                // Add an arc for each element, with a title for tooltips.
                const format = d3.format(",d");
                svg.append("g")
                    .attr("fill-opacity", 0.6)
                    .selectAll("path")
                    .data(root.descendants().filter(d => d.depth))
                    .join("path")
                    .attr("fill", (d) => { return color(d.data); })
                    .attr("d", arc)
                    .append("title")
                    .text(d => `${d.ancestors().map(d => d.data.name).reverse().join("/")}\n${format(d.value)}`);

                // Add a label for each element.
                svg.append("g")
                    .attr("pointer-events", "none")
                    .attr("text-anchor", "middle")
                    .attr("font-size", 10)
                    .attr("font-family", "sans-serif")
                    .selectAll("text")
                    .data(root.descendants().filter(d => d.depth && (d.y0 + d.y1) / 2 * (d.x1 - d.x0) > 10))
                    .join("text")
                    .attr("transform", function (d) {
                        const x = (d.x0 + d.x1) / 2 * 180 / Math.PI;
                        const y = (d.y0 + d.y1) / 2;
                        return `rotate(${x - 90}) translate(${y},0) rotate(${x < 180 ? 0 : 180})`;
                    })
                    .attr("dy", "0.35em")
                    .text(d => d.data.name);

                return () => {
                    if (d3Container?.current?.innerHTML) { d3Container.current.innerHTML = '' }
                }
            }
        }, [animal, d3Container.current])
    return <svg
        className="d3-component"
        width="100%"
        height="100%"
        ref={d3Container}
    />
}

export default AnimalRadialTree
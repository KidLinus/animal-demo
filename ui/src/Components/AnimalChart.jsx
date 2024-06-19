import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

const AnimalChart = ({ data = [] }) => {
    const d3Container = useRef(null);
    useEffect(
        () => {
            if (data && d3Container.current) {
                // const svg = d3.select(d3Container.current);
                // Specify the dimensions of the chart.
                const width = d3Container.current.getBoundingClientRect().width;
                const height = d3Container.current.getBoundingClientRect().height;

                // Specify the color scale.
                // const color = d3.scaleOrdinal(d3.schemeCategory10);

                // The force simulation mutates links and nodes, so create a copy
                // so that re-evaluating this cell produces the same result.
                // const links = data.links.map(d => ({ ...d }));
                const links = data.reduce((s, v) => [...s, ...Object.values(v.parents).map(id => ({ source: id, target: v.id, value: 1 }))], [])
                    .filter(({ source, target }) => data.find(d => d.id == source) && data.find(d => d.id == target))
                const nodes = data

                // Create a simulation with several forces.
                const simulation = d3.forceSimulation(nodes)
                    .force("link", d3.forceLink(links).id(d => d.id))
                    .force("charge", d3.forceManyBody())
                    .force("center", d3.forceCenter(width / 2, height / 2))
                    .on("tick", ticked);

                // Create the SVG container.
                const svg = d3.select(d3Container.current)
                    .attr("width", width)
                    .attr("height", height)
                    .attr("viewBox", [0, 0, width, height])
                    .attr("style", "max-width: 100%; height: auto;");

                // Add a line for each link, and a circle for each node.
                const link = svg.append("g")
                    .attr("stroke", "#999")
                    .attr("stroke-opacity", 0.6)
                    .selectAll()
                    .data(links)
                    .join("line")
                    .attr("stroke-width", d => Math.sqrt(d.value));

                const node = svg.append("g")
                    .selectAll()
                    .data(nodes)
                    .join("g")

                node.append("circle")
                    .attr("r", 5)
                    .attr("fill", d => d.gender == 0 ? "#00f" : "#F00");

                node.append("text")
                    .attr("x", 8)
                    .attr("y", 0)
                    .attr("font-size", "8px")
                    .text(d => `${d.id} ${d.name}`)
                
                    // Add a drag behavior.
                node.call(d3.drag()
                    .on("start", dragstarted)
                    .on("drag", dragged)
                    .on("end", dragended));

                function ticked() {
                    link
                        .attr("x1", d => d.source.x)
                        .attr("y1", d => d.source.y)
                        .attr("x2", d => d.target.x)
                        .attr("y2", d => d.target.y);

                    node.attr("transform", d => "translate(" + [d.x, d.y] + ")")
                }

                function dragstarted(event) {
                    if (!event.active) simulation.alphaTarget(0.3).restart();
                    event.subject.fx = event.subject.x;
                    event.subject.fy = event.subject.y;
                }

                // Update the subject (dragged node) position during drag.
                function dragged(event) {
                    event.subject.fx = event.x;
                    event.subject.fy = event.y;
                }

                // Restore the target alpha so the simulation cools after dragging ends.
                // Unfix the subject position now that it’s no longer being dragged.
                function dragended(event) {
                    if (!event.active) simulation.alphaTarget(0);
                    event.subject.fx = null;
                    event.subject.fy = null;
                }
                return () => {
                    simulation.stop()
                    if (d3Container?.current?.innerHTML) { d3Container.current.innerHTML = '' }
                }
            }
        }, [data, d3Container.current])
    return <svg
        className="d3-component"
        width="100%"
        height="100%"
        ref={d3Container}
    />
}

export default AnimalChart
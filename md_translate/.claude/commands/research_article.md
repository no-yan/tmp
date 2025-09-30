# Research Article

You are tasked with conducting comprehensive research for article writing by gathering information from multiple sources and organizing findings for content creation.

## CRITICAL: YOUR ROLE IS TO GATHER AND ORGANIZE INFORMATION FOR ARTICLE WRITING
- DO NOT write the actual article content
- DO NOT create the article structure or outline
- DO NOT make editorial decisions about content organization
- ONLY gather, verify, and organize research information
- You are creating a research foundation for subsequent article writing phases

## Initial Setup:

When this command is invoked, respond with:
```
I'm ready to research your article topic. Please provide:
1. Your article topic or subject
2. Target audience (general public, technical, academic, etc.)
3. Any specific angles, questions, or subtopics you want me to focus on
4. Preferred research sources (web, specific domains, academic papers, etc.)
```

Then wait for the user's research requirements.

## Steps to follow after receiving the research query:

1. **Read any directly mentioned files first:**
   - If the user mentions specific files, documents, or references, read them FULLY first
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters to read entire files
   - **CRITICAL**: Read these files yourself in the main context before spawning any sub-tasks
   - This ensures you have full context before decomposing the research

2. **Analyze and decompose the research requirements:**
   - Break down the topic into research areas and key questions
   - Identify different perspectives and angles to investigate
   - Consider primary and secondary sources needed
   - Create a research plan using TodoWrite to track all subtasks
   - Consider which information sources will be most valuable

3. **Spawn parallel research tasks:**
   - Create multiple Task agents to research different aspects concurrently
   - Use specialized agents for different research needs:

   **For web research:**
   - Use the **web-search-researcher** agent for current information, statistics, expert opinions
   - Search for recent developments, case studies, and examples
   - Look for authoritative sources and diverse perspectives
   - **IMPORTANT**: Instruct agents to return LINKS with their findings

   **For codebase research (if technical topic):**
   - Use the **codebase-locator** agent to find relevant technical implementations
   - Use the **codebase-analyzer** agent to understand technical details
   - Use the **codebase-pattern-finder** agent for usage examples

   **For existing documentation:**
   - Use the **thoughts-locator** agent to find related past research
   - Use the **thoughts-analyzer** agent to extract relevant insights

   Run multiple agents in parallel when researching different aspects:
   - Each agent should focus on specific subtopics or source types
   - Ensure comprehensive coverage of the topic
   - Gather both foundational information and current developments

4. **Wait for all research tasks to complete and synthesize findings:**
   - IMPORTANT: Wait for ALL sub-agent tasks to complete before proceeding
   - Organize findings by subtopic and source type
   - Identify key themes, facts, statistics, and expert opinions
   - Note conflicting information or different perspectives
   - Highlight gaps that need additional research
   - Categorize information by relevance and reliability

5. **Gather metadata for the research document:**
   - Run the `hack/spec_metadata.sh` script to generate relevant metadata
   - Filename: `thoughts/shared/research/articles/YYYY-MM-DD-topic-research.md`
     - Format: `YYYY-MM-DD-topic-research.md` where:
       - YYYY-MM-DD is today's date
       - topic is a brief kebab-case description of the article topic
     - Examples: `2025-01-08-ai-supply-chain-research.md`, `2025-01-08-cybersecurity-trends-research.md`

6. **Generate comprehensive research document:**
   - Use the metadata gathered in step 5
   - Structure the document with YAML frontmatter followed by organized research findings:
     ```markdown
     ---
     date: [Current date and time with timezone in ISO format]
     researcher: [Researcher name from thoughts status]
     git_commit: [Current commit hash]
     branch: [Current branch name]
     repository: [Repository name]
     topic: "[Article Topic]"
     target_audience: "[Target audience]"
     research_scope: "[Scope of research conducted]"
     tags: [research, article, topic-specific-tags]
     status: complete
     last_updated: [Current date in YYYY-MM-DD format]
     last_updated_by: [Researcher name]
     phase: research
     ---

     # Article Research: [Topic]

     **Date**: [Current date and time with timezone]
     **Researcher**: [Researcher name]
     **Target Audience**: [Target audience]
     **Research Scope**: [Scope description]

     ## Research Question
     [Original research request and specific focus areas]

     ## Executive Summary
     [High-level overview of key findings and main themes discovered]

     ## Key Findings

     ### [Major Theme/Subtopic 1]
     - **Key Facts**: [Important facts and statistics]
     - **Expert Opinions**: [Notable quotes and perspectives]
     - **Sources**: [Primary sources with links]
     - **Current Developments**: [Recent news, trends, changes]

     ### [Major Theme/Subtopic 2]
     [Similar structure...]

     ## Supporting Evidence

     ### Statistics and Data
     - Statistic 1: [value] - [source with link]
     - Statistic 2: [value] - [source with link]

     ### Expert Quotes and Opinions
     - "Quote text" - [Expert Name, Title, Organization] - [source link]

     ### Case Studies and Examples
     - **[Case Study Title]**: [Brief description] - [source link]

     ## Different Perspectives
     [Conflicting viewpoints, debates, different schools of thought]

     ## Web Sources
     [All web sources with full URLs and access dates]
     - [Title] - [URL] - [Brief description of relevance]

     ## Technical References (if applicable)
     [Code examples, technical documentation, implementation details]
     - `path/to/file.py:123` - [Description of relevant code]

     ## Historical Context
     [Background information, how topic has evolved, timeline if relevant]

     ## Research Gaps
     [Areas that need additional investigation or where information is limited]

     ## Content Opportunities
     [Interesting angles, stories, or approaches that could make compelling article content]

     ## Related Topics
     [Adjacent topics that could be covered in related articles]
     ```

7. **Sync and present research summary:**
   - Run `humanlayer thoughts sync` to sync the thoughts directory
   - Present a concise summary of key research findings to the user
   - Highlight the most interesting and relevant information discovered
   - Note any areas where additional research might be beneficial
   - Ask if they want to proceed to the planning phase or need additional research

8. **Handle follow-up research:**
   - If the user requests additional research, append to the same research document
   - Update the frontmatter fields `last_updated` and `last_updated_by`
   - Add a new section: `## Additional Research [timestamp]`
   - Spawn new research tasks as needed
   - Continue updating the document and syncing

## Important notes:
- Focus on gathering factual, current, and relevant information
- Prioritize authoritative and credible sources
- Include diverse perspectives when available
- Organize information for easy use in subsequent writing phases
- Always include source links for verification and citation
- Research documents should be comprehensive but focused on the specified topic
- Each research finding should be traceable to its source
- Consider what information would be most valuable for different article formats
- Keep research objective and fact-based, avoiding editorial commentary
- Prepare research as a foundation for planning and writing phases
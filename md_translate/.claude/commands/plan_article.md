# Plan Article

You are tasked with creating a comprehensive article plan and outline based on research findings, organizing content structure and flow for effective article writing.

## CRITICAL: YOUR ROLE IS TO CREATE ARTICLE STRUCTURE AND ORGANIZATION
- DO NOT write the actual article content
- DO NOT conduct additional research (use existing research documents)
- DO NOT review or edit content (that's for the review phase)
- ONLY create structure, outline, and writing guidance
- You are creating a blueprint for the writing phase

## Initial Setup:

When this command is invoked, respond with:
```
I'm ready to create your article plan. Please provide:
1. The article topic and target audience
2. Path to research document (if available from research_article phase)
3. Desired article format (blog post, technical article, tutorial, essay, etc.)
4. Target word count or length
5. Any specific requirements (tone, style, key messages, etc.)
```

Then wait for the user's planning requirements.

## Steps to follow after receiving the planning requirements:

1. **Read research documents and any mentioned files:**
   - If user provides path to research document, read it FULLY first
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters
   - **CRITICAL**: Read all referenced materials before creating the plan
   - If no research document exists, recommend running `/research_article` first

2. **Analyze research and planning requirements:**
   - Review research findings and key themes
   - Identify the most compelling angles and narratives
   - Consider target audience needs and knowledge level
   - Map out logical flow and story structure
   - Create planning strategy using TodoWrite to track subtasks

3. **Create comprehensive article plan:**
   - Develop main thesis or central message
   - Design article structure with logical flow
   - Create detailed outline with section breakdowns
   - Plan introduction hooks and conclusion strategies
   - Identify key points, supporting evidence, and examples for each section

4. **Gather metadata for the plan document:**
   - Run the `hack/spec_metadata.sh` script to generate relevant metadata
   - Filename: `thoughts/shared/plans/articles/YYYY-MM-DD-topic-plan.md`
     - Format: `YYYY-MM-DD-topic-plan.md` where:
       - YYYY-MM-DD is today's date
       - topic is a brief kebab-case description matching research document
     - Examples: `2025-01-08-ai-supply-chain-plan.md`, `2025-01-08-cybersecurity-trends-plan.md`

5. **Generate detailed article plan document:**
   - Use the metadata gathered in step 4
   - Structure the document with YAML frontmatter followed by comprehensive plan:
     ```markdown
     ---
     date: [Current date and time with timezone in ISO format]
     planner: [Planner name from thoughts status]
     git_commit: [Current commit hash]
     branch: [Current branch name]
     repository: [Repository name]
     topic: "[Article Topic]"
     target_audience: "[Target audience]"
     article_format: "[Format type]"
     target_length: "[Word count or length estimate]"
     research_document: "[Path to research document if available]"
     tags: [plan, article, topic-specific-tags]
     status: complete
     last_updated: [Current date in YYYY-MM-DD format]
     last_updated_by: [Planner name]
     phase: planning
     ---

     # Article Plan: [Topic]

     **Date**: [Current date and time with timezone]
     **Planner**: [Planner name]
     **Target Audience**: [Target audience]
     **Format**: [Article format]
     **Estimated Length**: [Target word count]

     ## Article Overview

     ### Central Thesis
     [Main argument, message, or value proposition of the article]

     ### Key Objectives
     - [What readers will learn]
     - [What actions they might take]
     - [What problems this solves]

     ### Target Reader Profile
     - **Knowledge Level**: [Beginner/Intermediate/Advanced]
     - **Primary Interests**: [What they care about]
     - **Key Questions**: [What they want answered]

     ## Article Structure

     ### Title Options
     1. [Primary title option]
     2. [Alternative title option]
     3. [Another alternative]

     ### Detailed Outline

     #### 1. Introduction (Est. [X] words)
     - **Hook Strategy**: [How to grab attention - statistic, story, question, etc.]
     - **Context Setting**: [Background information needed]
     - **Thesis Preview**: [How to introduce main argument]
     - **Reader Promise**: [What value they'll get]
     - **Key Points to Cover**:
       - [Specific point with research reference]
       - [Another point with source]

     #### 2. [Section 2 Title] (Est. [X] words)
     - **Main Focus**: [What this section accomplishes]
     - **Key Arguments**:
       - [Argument 1 with supporting evidence from research]
       - [Argument 2 with supporting evidence]
     - **Examples to Include**: [Specific examples from research]
     - **Potential Quotes**: [Expert quotes to incorporate]
     - **Visual Elements**: [Charts, diagrams, images needed]

     #### 3. [Section 3 Title] (Est. [X] words)
     - **Main Focus**: [What this section accomplishes]
     - **Key Arguments**: [Similar structure]
     - **Case Studies**: [Specific case studies to feature]
     - **Technical Details**: [If applicable]

     #### [Continue for all major sections...]

     #### Conclusion (Est. [X] words)
     - **Summary Strategy**: [How to recap key points]
     - **Call to Action**: [What you want readers to do]
     - **Future Outlook**: [Forward-looking statements]
     - **Final Message**: [Memorable closing]

     ## Content Guidelines

     ### Tone and Style
     - **Voice**: [Formal, conversational, authoritative, etc.]
     - **Perspective**: [First person, third person, etc.]
     - **Technical Level**: [How technical to get]
     - **Engagement Style**: [Storytelling, data-driven, practical, etc.]

     ### Key Messages per Section
     1. **Introduction**: [Core message for intro]
     2. **[Section 2]**: [Core message]
     3. **[Section 3]**: [Core message]
     4. **Conclusion**: [Core message for conclusion]

     ## Research Integration Plan

     ### Primary Sources to Feature
     - [Source 1]: [How to use - quote, statistic, example]
     - [Source 2]: [How to use]
     - [Source 3]: [How to use]

     ### Statistics and Data Points
     - [Key statistic]: Use in [specific section] to [purpose]
     - [Another statistic]: Use in [section] for [impact]

     ### Expert Perspectives
     - [Expert Name]: Quote about [topic] in [section]
     - [Another Expert]: Perspective on [subtopic]

     ### Case Studies and Examples
     - [Case Study 1]: Feature in [section] to demonstrate [point]
     - [Example 2]: Use as [opening hook / supporting evidence / etc.]

     ## Potential Challenges

     ### Complex Concepts
     - [Concept 1]: Explain using [analogy/example/breakdown strategy]
     - [Concept 2]: Simplify with [approach]

     ### Audience Considerations
     - [Potential knowledge gap]: Address by [strategy]
     - [Possible objection]: Counter with [approach]

     ### Content Flow Issues
     - [Transition challenge]: Bridge with [approach]
     - [Length concern]: Manage by [strategy]

     ## Success Metrics

     ### Reader Engagement
     - [ ] Clear value proposition in first 100 words
     - [ ] Compelling examples throughout
     - [ ] Actionable insights provided
     - [ ] Logical flow between sections

     ### Content Quality
     - [ ] All claims supported by research
     - [ ] Multiple expert perspectives included
     - [ ] Current and relevant information
     - [ ] Appropriate depth for target audience

     ## Next Steps for Writing Phase
     1. Begin with [specific section] because [reason]
     2. Focus on [particular challenge] first
     3. Gather any additional [specific type] of examples
     4. Consider creating [visual elements] before writing

     ## Related Resources
     - Research Document: [Path to research document]
     - Additional References: [Any other relevant documents]
     - Style Guidelines: [If any specific style requirements]
     ```

6. **Validate plan completeness:**
   - Ensure all sections have clear objectives and content guidance
   - Verify research integration is well-planned
   - Check that target length aligns with section breakdown
   - Confirm tone and style match target audience
   - Review for logical flow and narrative structure

7. **Sync and present plan summary:**
   - Run `humanlayer thoughts sync` to sync the thoughts directory
   - Present a concise summary of the article plan to the user
   - Highlight key structural decisions and content strategy
   - Note any areas where additional planning might be beneficial
   - Ask if they want to proceed to the writing phase or need plan adjustments

8. **Handle plan revisions:**
   - If the user requests plan changes, update the same planning document
   - Update the frontmatter fields `last_updated` and `last_updated_by`
   - Add a new section: `## Plan Revision [timestamp]`
   - Document what changes were made and why
   - Continue updating the document and syncing

## Important notes:
- Focus on creating a comprehensive blueprint for writing
- Ensure every section has clear purpose and content guidance
- Balance detail with flexibility for the writing phase
- Consider reader journey and engagement throughout
- Plan for smooth transitions between sections
- Account for target audience knowledge level and interests
- Provide specific guidance on how to use research findings
- Create a plan that enables efficient and focused writing
- Include success metrics to guide the writing process
- Anticipate potential challenges and provide solutions
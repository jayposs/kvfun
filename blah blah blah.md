# Why KVFun

There were 2 driving forces for creation of this code. First I was working on a Javascript project that provided nice data handling features in the browser. Second I was exploring a fairly new database called HarperDB (HDB). The main attractions of it were simplicity and excellent features. I'm still a fan of HDB, but worried about reliability and complexity of setup and upgrades. To be fair, with Docker, the setup process is very easy. But I'm pretty sure what's actually going on with the setup is quite complicated. The foundation of HDB is LMDB (same logic as BoltDB), but it also implements a sql interface on top.

Between my JS project and HarperDB the idea of creating a KV frontend emerged.

Previously I used MongoDB and various Sql dbs. If using a service, like Mongo's Atlas, then things are really simple. Of course there may be a problems due to long distance connections and increased cost. 

With BoltDB and KVFun enhancements I get a simple setup and use solution. The feature set is definitely below other tools but for most transactional applications I've written, don't think that's a big problem. For jobs like data analysis, the data needs to be accessible by a more flexible tool (like Snowflake).

When setting up a system, there are a million things to consider. I just wanted to create a solution that focused on being simple, reliable, and cheap but had enough features to be  useful.

**Caveat -**  Anytime you're working with Go modules and Git/github there's a certain level of complexity and mystery setting things up. 

If you decide to explore kvfun, consider it your own code with which to modify as needed. I really tried to keep the code small and easy to understand. See notes.md for more information on how to add new features. I don't consider kvfun a full featured package, but rather a good starting point.  

I'm 66 years old (as of 2023) so my future plans are rather uncertain. Have Fun.  
Jay

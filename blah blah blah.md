# Why KVFun - Quest for Simplicity  

I recently explored a fairly new database called HarperDB (HDB). The main attractions were simplicity and excellent features. I'm still a fan of HDB, but worried about reliability and complexity of setup and upgrades. To be fair, with Docker, the setup process is very easy. But I'm pretty sure what's actually going on with the setup is quite complicated. The foundation of HDB is LMDB (same logic as BoltDB), but it also implements a sql interface on top.

Previously I used MongoDB and various Sql dbs. If using a service, like Mongo's Atlas, then things are really simple. Of course there may be a speed problem due to long distance connections and increased cost. 

With BoltDB and KVFun enhancements I get a simple setup and use solution. The feature set is definitely below other tools but for most transactional applications I've written, don't think that's a big problem. When it comes to reporting and data analysis, the data needs to be accessible by a more flexible tool (like Snowflake).

When setting up a system, there are a million things to consider. I just wanted an option that focused on being simple and cheap but had enough features to be really useful.

**Caveat**  Anytime you're working with Go modules and Git/github there's a certain level of complexity and mystery setting things up. 
